package pc

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	ap "pc/argparse"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	pluralize "github.com/gertd/go-pluralize"
	"golang.org/x/exp/slices"
)

type T_rawdata []string
type T_dataline []string
type T_parsedData []T_dataline

func printJSON(d T_parsedData) {
	//  set seperator rune
	sep := []rune(ap.CmdParams.Sep)[0]
	// split headerline
	header := LineParse(ap.CmdParams.Header, sep)
	fmt.Println("[")
	lines := len(d)
	for ln, line := range d {
		fmt.Println("  {")
		l := len(line)
		for col, val := range line {
			fmt.Print("    \"", header[col], "\": \"", val, "\"")
			if col+1 < l {
				fmt.Println(",")
			} else {
				fmt.Println("")
			}
		}
		if ln+1 < lines {
			fmt.Println("  },")
		} else {
			fmt.Println("  }")
		}
	}
	fmt.Println("]")
}

func printJSONwithTC(d T_parsedData) {
	//  set seperator rune
	sep := []rune(ap.CmdParams.Sep)[0]
	// split headerline
	header := d[0]
	if ap.CmdParams.Header != "" {
		header = LineParse(ap.CmdParams.Header, sep)
	}
	if ap.CmdParams.Ts {
		d = d[1:]
	}
	hkey := header[0]
	header = header[1:]
	lines := len(d)

	pl := pluralize.NewClient()
	fmt.Println("{")
	fmt.Print("  \"", pl.Plural(hkey), "\": [")
	fmt.Println("")
	for ln, line := range d {
		fmt.Println("    {")
		hval := line[0]
		fmt.Print("      \"", hkey, "\": \"", hval, "\"")
		fmt.Println(",")
		fmt.Println("      \"data\": {")
		line = line[1:]
		l := len(line)
		for col, val := range line {
			fmt.Print("        \"", header[col], "\": \"", val, "\"")
			if col+1 < l {
				fmt.Println(",")
			} else {
				fmt.Println("")
			}
		}
		fmt.Println("      }")
		if ln+1 < lines {
			fmt.Println("    },")
		} else {
			fmt.Println("    }")
		}
	}
	fmt.Println("  ]")
	fmt.Println("}")
}

// PrintJson prints data as JSON
func (d T_parsedData) PrintJson() {
	if ap.CmdParams.Header == "" && !ap.CmdParams.Ts {
		b, err := json.MarshalIndent(d, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		if ap.CmdParams.Jtc || ap.CmdParams.Ts {
			// use first column as key
			printJSONwithTC(d)
		} else {
			printJSON(d)
		}
	}
}

// PrintCsv prints data as CSV
func (d T_parsedData) PrintCsv() {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	for _, record := range d {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record", err)
		}
	}
}

// Append appends a dataline
func (d *T_parsedData) Append(l T_dataline) {
	*d = append(*d, l)
	// return d
}

// Insert inserts a dataline at given position
func (d *T_parsedData) Insert(l T_dataline, pos int) {
	if pos >= 0 && pos <= len(*d) {
		*d = append(*d, l)
		data := *d
		copy(data[(pos+1):], data[pos:])
		data[pos] = l
		d = &data
	}
}

// generateLine format the fields of dataline to maxlen of columns
func (data *T_dataline) generateLine(maxlen T_maxlenghts) {
	// fmt.Println("dataline:", *data)
	for pos, mxlen := range maxlen {
		val := ""
		if pos < len((*data)) {
			val = (*data)[pos]
			runecount := utf8.RuneCountInString(val)
			blanklen := mxlen - runecount
			// fmt.Println("pos:", pos, "mxlen:", mxlen, "runecount:", runecount, "blanklen:", blanklen, "val:", val)
			if regexp.MustCompile(`^ *[0-9\.,]+ *$`).MatchString((*data)[pos]) && !ap.CmdParams.Nn {
				(*data)[pos] = strings.Repeat(" ", blanklen) + val
			} else {
				(*data)[pos] = val + strings.Repeat(" ", blanklen)
			}
		} else {
			(*data) = append((*data), strings.Repeat(" ", mxlen))
		}
	}
}

// InsertGroupSeperator inserts a trenenr slice when the content of gcol changed the value.
// Let further values of gcol empty until the next group change
func (data *T_parsedData) InsertGroupSeperator(gcol int, gcolval bool, trenner []string) {
	nd := T_parsedData{}
	if gcol > 0 && gcol <= len(*data)+1 {
		gcol -= 1
		ref := ""
		if len((*data)[0]) > gcol {
			ref = (*data)[0][gcol]
		}
		for i, row := range *data {
			if len(row) > gcol && ref != row[gcol] && ref != trenner[gcol] && row[gcol] != trenner[gcol] && i > 0 {
				nd.Append(trenner)
				ref = row[gcol]
				nd.Append(row)
			} else {
				if !gcolval && i > 0 && len(row) > gcol && ref == row[gcol] && row[gcol] != trenner[gcol] {
					row[gcol] = "''"
				}
				nd.Append(row)
				if ref == trenner[gcol] {
					ref = row[gcol]
				}
			}
		}
		*data = nd
	}
}

// selectColumns select data columns as defined in CmdParms.Columns
func (data *T_parsedData) selectColumns() {
	if len(ap.CmdParams.Columns) > 0 {
		for i, row := range *data {
			nrow := T_dataline{}
			for _, col := range ap.CmdParams.Columns {
				if int(col) > 0 && int(col) <= len(row) {
					nrow = append(nrow, row[col-1])
				} else {
					nrow = append(nrow, "")
				}
			}
			(*data)[i] = nrow
		}
	}
}

// selectColumns select columns from dataline as defined in CmdParms.Columns
func (data *T_dataline) selectColumns() {
	nrow := T_dataline{}
	for _, col := range ap.CmdParams.Columns {
		if int(col) > 0 && int(col) <= len(*data) {
			nrow = append(nrow, (*data)[col-1])
		} else {
			nrow = append(nrow, "")
		}
	}
	*data = nrow
}

// insertTrenner insert trenner for TitelSeperator, FooterSeperator or PrettyPrint
func (data *T_parsedData) insertTrenner(trenner []string) {
	if ap.CmdParams.Ts || ap.CmdParams.Fs || ap.CmdParams.Pp {
		if ap.CmdParams.Pp {
			if ap.CmdParams.Fs {
				(*data).Insert(trenner, len(*data)-1)
			}
			(*data).Insert(trenner, 0)
			(*data).Insert(trenner, 2)
			(*data).Append(trenner)
		} else {
			if ap.CmdParams.Ts {
				(*data).Insert(trenner, 1)
			}
			if ap.CmdParams.Fs {
				(*data).Insert(trenner, len((*data))-1)
			}
		}
	}
}

// sort data on k column
func (data *T_parsedData) sort(k int) {
	k -= 1
	l1 := T_dataline{}
	d := T_parsedData{}
	if !ap.CmdParams.Nhl {
		l1 = (*data)[0]
		d = (*data)[1:]
	} else {
		d = *data
	}
	sort.SliceStable(d, func(i, j int) bool {
		return d[i][k] < d[j][k]
	})
	da := T_parsedData{}
	if !ap.CmdParams.Nhl {
		da = append(da, l1)
	}
	da = append(da, d...)
	*data = da
}

// delete elements from data
func (data *T_parsedData) delete(i, j int) {
	*data = slices.Delete(*data, i, j)
}

// format data to column max width
func (data *T_parsedData) formatDataToMaxWidth(maxlen []int) {
	for i := range *data {
		(*data)[i].generateLine(maxlen)
	}
}

// printAsciiTab prints the data as ASCII table
func (data *T_parsedData) printAsciiTab() {
	sp := strings.Repeat(" ", ap.CmdParams.ColSepW)
	for _, row := range *data {
		var line string
		if ap.CmdParams.Pp || ap.CmdParams.Cs {
			line = ap.CmdParams.Colsep + sp + strings.Join(row, sp+ap.CmdParams.Colsep+sp) + sp + ap.CmdParams.Colsep
		} else {
			line = strings.Join(row, sp)
		}
		if ap.CmdParams.MoreBlanks {
			line = strings.Replace(line, "§", " ", -1)
		}
		fmt.Println(line)
	}
}

// Format select the data and header columns if CmdParams.Colums is set.
// Insert trenner for Ts, Fs and Pp. Format the data fields to maxlen for each column.
// Print the data as CSV, JSON or ASCII table depending on options.
func Format(data T_parsedData) {
	//  set seperator rune
	sep := []rune(ap.CmdParams.Sep)[0]

	// select data columns
	if len(ap.CmdParams.Columns) > 0 {
		data.selectColumns()
	}

	// remove header line, first row
	if ap.CmdParams.Rh {
		data.delete(0, 1)
	}

	if ap.CmdParams.SortCol > 0 {
		data.sort(int(ap.CmdParams.SortCol))
	}

	// Insert Headerline from CmdParams
	if ap.CmdParams.Header != "" && !ap.CmdParams.Json {
		headerline := LineParse(ap.CmdParams.Header, sep)
		if len(ap.CmdParams.Columns) > 0 && len(headerline) > len(ap.CmdParams.Columns) {
			// select columns from header
			headerline.selectColumns()
		}
		data.Insert(headerline, 0)
	}

	// get maxlen for each column
	maxlen := GetMaxLength(data)

	// insert a row with col numbers
	if ap.CmdParams.Num {
		n := []string{}
		for i := range maxlen {
			ns := strconv.Itoa(i + 1)
			if len(ap.CmdParams.Columns) > 0 {
				ns = ns + " [" + strconv.Itoa(int(ap.CmdParams.Columns[i])) + "]"
			}
			n = append(n, ns)
		}
		data.Insert(n, 0)
		// get maxlen for each column, must calculated again aafter inserting the number line
		maxlen = GetMaxLength(data)
	}

	//  define trenner slice
	trenner := []string{}
	for _, v := range maxlen {
		trenner = append(trenner, strings.Repeat("-", v))
	}

	// insert trenner for TitelSeperator, FooterSeperator or PrettyPrint
	if !(ap.CmdParams.Json || ap.CmdParams.Csv) {
		data.insertTrenner(trenner)
	}

	// print CSV and return
	if ap.CmdParams.Csv {
		data.PrintCsv()
		return
	}

	// print JSON and return
	if ap.CmdParams.Json {
		data.PrintJson()
		return
	}

	// insert trenner between GroupChange of gcol
	data.InsertGroupSeperator(int(ap.CmdParams.Gcol), ap.CmdParams.GcolVal, trenner)

	// format all columns for maxlen column width
	if !ap.CmdParams.Nf {
		data.formatDataToMaxWidth(maxlen)
	}

	// print data slices as line with or without column seperator
	data.printAsciiTab()
}
