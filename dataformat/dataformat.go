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

// printJSON prints the parsed data in JSON format.
// It uses the header line defined in CmdParams.Header and the separator defined in CmdParams.Sep.
// The function creates a JSON array where each element is an object representing a row of data.
// Each object contains key-value pairs, where the keys are the column headers and the values are the corresponding data fields.
func printJSON(d T_parsedData) {
	// Set separator rune
	sep := []rune(ap.CmdParams.Sep)[0]
	// Split header line using the separator
	header := LineParse(ap.CmdParams.Header, sep)
	// Print the opening bracket of the JSON array
	fmt.Println("[")
	lines := len(d)
	// Iterate over each line of data
	for ln, line := range d {
		// Print the opening brace of the JSON object
		fmt.Println("  {")
		l := len(line)
		// Iterate over each column value in the line
		for col, val := range line {
			// Print the key-value pair for the current column
			fmt.Print("    \"", header[col], "\": \"", val, "\"")
			// If not the last column, print a comma
			if col+1 < l {
				fmt.Println(",")
			} else {
				fmt.Println("")
			}
		}
		// Print the closing brace of the JSON object
		if ln+1 < lines {
			fmt.Println("  },")
		} else {
			fmt.Println("  }")
		}
	}
	// Print the closing bracket of the JSON array
	fmt.Println("]")
}

// printJSONwithTC prints the parsed data in JSON format with a top-level collection.
// It uses the header line defined in CmdParams.Header and the separator defined in CmdParams.Sep.
// If the TitleSeparator (Ts) flag is set, it skips the first line of data.
// The function creates a JSON object where the top-level key is the plural form of the first header field,
// and the value is an array of objects. Each object represents a row of data, with the first column as a key-value pair
// and the remaining columns nested under a "data" key.
func printJSONwithTC(d T_parsedData) {
	// Set separator rune
	sep := []rune(ap.CmdParams.Sep)[0]
	// Split header line
	header := d[0]
	if ap.CmdParams.Header != "" {
		header = LineParse(ap.CmdParams.Header, sep)
	}
	// If TitleSeparator flag is set, skip the first line of data
	if ap.CmdParams.Ts {
		d = d[1:]
	}
	// Extract the first header key and the remaining headers
	hkey := header[0]
	header = header[1:]
	lines := len(d)
	// Initialize pluralize client
	pl := pluralize.NewClient()
	// Print the opening of the JSON object
	fmt.Println("{")
	fmt.Print("  \"", pl.Plural(hkey), "\": [")
	fmt.Println("")
	// Iterate over each line of data
	for ln, line := range d {
		fmt.Println("    {")
		// Extract the value for the first header key
		hval := line[0]
		fmt.Print("      \"", hkey, "\": \"", hval, "\"")
		fmt.Println(",")
		fmt.Println("      \"data\": {")
		// Extract the remaining line values
		line = line[1:]
		l := len(line)
		// Iterate over each column value and print it as a key-value pair
		for col, val := range line {
			fmt.Print("        \"", header[col], "\": \"", val, "\"")
			if col+1 < l {
				fmt.Println(",")
			} else {
				fmt.Println("")
			}
		}
		fmt.Println("      }")
		// Print the closing of the current object
		if ln+1 < lines {
			fmt.Println("    },")
		} else {
			fmt.Println("    }")
		}
	}
	// Print the closing of the JSON array and object
	fmt.Println("  ]")
	fmt.Println("}")
}

// PrintJson prints the parsed data in JSON format.
// If no header is specified and the TitleSeparator (Ts) flag is not set, it marshals the data directly to JSON.
// Otherwise, it uses either printJSONwithTC or printJSON based on the CmdParams flags.
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

// PrintCsv prints the parsed data in CSV format.
// It writes the data to the standard output using the csv.Writer from the encoding/csv package.
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

// InsertGroupSeperator inserts a separator slice when the content of the specified column changes value.
// It takes a column index, a boolean flag, and two slices (trenner and htrenner) as input.
// The function inserts a separator slice whenever the value in the specified column changes.
// It leaves further values of the column empty until the next group change.
func (data *T_parsedData) InsertGroupSeperator(gcol int, gcolval bool, trenner, htrenner []string) {
	nd := T_parsedData{}
	if gcol > 0 && gcol <= len(*data)+1 {
		gcol -= 1
		ref := ""
		if len((*data)[0]) > gcol {
			ref = (*data)[0][gcol]
		}
		for i, row := range *data {
			// if we not on first line and not on a trenner line
			if i > 0 && len(row) > gcol && row[gcol] != trenner[gcol] && row[gcol] != htrenner[gcol] {
				// when ref (thats the field of the previous line) differs from the filed of the current line
				if ref != row[gcol] && ref != trenner[gcol] && ref != htrenner[gcol] {
					nd.Append(trenner)
					//  set ref to the new field value
					ref = row[gcol]
					nd.Append(row)
				} else {
					//  if no difference in the field and gcolval is not set
					if !gcolval && ref == row[gcol] {
						row[gcol] = "''"
					}
					nd.Append(row)
				}
			} else {
				nd.Append(row)
				// ref = row[gcol]
			}
		}
		*data = nd
	}
}

// selectColumns selects data columns as defined in CmdParams.Columns.
// It iterates over each row in the parsed data and creates a new row containing only the selected columns.
// If a column index is out of range, it adds an empty string to the new row.
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

// selectColumns selects columns from dataline as defined in CmdParams.Columns.
// It creates a new dataline containing only the selected columns.
// If a column index is out of range, it adds an empty string to the new dataline.
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

// insertTrenner inserts separators for TitleSeparator, FooterSeparator, or PrettyPrint.
// It takes two slices (trenner and htrenner) as input and inserts them at appropriate positions
// based on the CmdParams flags.
func (data *T_parsedData) insertTrenner(trenner, htrenner []string) {
	if ap.CmdParams.Ts || ap.CmdParams.Fs || ap.CmdParams.Pp {
		if ap.CmdParams.Pp {
			if ap.CmdParams.Fs {
				(*data).Insert(htrenner, len(*data)-1)
			}
			(*data).Insert(trenner, 0)
			(*data).Insert(trenner, 2)
			(*data).Append(trenner)
		} else {
			if ap.CmdParams.Ts {
				(*data).Insert(htrenner, 1)
			}
			if ap.CmdParams.Fs {
				(*data).Insert(htrenner, len((*data))-1)
			}
		}
	}
}

// sort sorts the parsed data based on the specified column index.
// It takes an integer k as input, which represents the column index to sort by (1-based index).
// If the NoHeaderLine (Nhl) flag is not set, it treats the first row as a header and excludes it from sorting.
// The function uses a stable sort to maintain the relative order of rows with equal values in the specified column.
func (data *T_parsedData) sort(k int) {
	// Convert 1-based index to 0-based index
	k -= 1
	// Initialize variables for header line and data to be sorted
	l1 := T_dataline{}
	d := T_parsedData{}
	// If NoHeaderLine flag is not set, separate the header line from the data
	if !ap.CmdParams.Nhl {
		l1 = (*data)[0]
		d = (*data)[1:]
	} else {
		// Otherwise, include all data for sorting
		d = *data
	}
	// Perform a stable sort on the data based on the specified column
	sort.SliceStable(d, func(i, j int) bool {
		return d[i][k] < d[j][k]
	})
	// Initialize a new parsed data slice to store the sorted data
	da := T_parsedData{}
	// If NoHeaderLine flag is not set, add the header line back to the sorted data
	if !ap.CmdParams.Nhl {
		da = append(da, l1)
	}
	// Append the sorted data to the new parsed data slice
	da = append(da, d...)
	// Update the original data with the sorted data
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

// printAsciiTab prints the parsed data as an ASCII table.
// It formats each row of the data according to the specified column separator and column width.
// If the PrettyPrint (Pp) or ColumnSeparator (Cs) flags are set, it adds the column separator between columns.
// If the MoreBlanks flag is set, it replaces the placeholder character '§' with spaces.
func (data *T_parsedData) printAsciiTab() {
	// Create a string of spaces with the length of the column separator width
	sp := strings.Repeat(" ", ap.CmdParams.ColSepW)
	// Iterate over each row in the parsed data
	for _, row := range *data {
		var line string
		// If PrettyPrint or ColumnSeparator flags are set, format the row with column separators
		if ap.CmdParams.Pp || ap.CmdParams.Cs {
			line = ap.CmdParams.Colsep + sp + strings.Join(row, sp+ap.CmdParams.Colsep+sp) + sp + ap.CmdParams.Colsep
		} else {
			// Otherwise, join the columns with spaces
			line = strings.Join(row, sp)
		}
		// If MoreBlanks flag is set, replace '§' with spaces
		if ap.CmdParams.MoreBlanks {
			line = strings.Replace(line, "§", " ", -1)
		}
		// Print the formatted line
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
	htrenner := []string{}
	for _, v := range maxlen {
		trenner = append(trenner, strings.Repeat("-", v))
		htrenner = append(htrenner, strings.Repeat("=", v))
	}
	// insert trenner for TitelSeperator, FooterSeperator or PrettyPrint
	if !(ap.CmdParams.Json || ap.CmdParams.Csv) {
		data.insertTrenner(trenner, htrenner)
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
	data.InsertGroupSeperator(int(ap.CmdParams.Gcol), ap.CmdParams.GcolVal, trenner, htrenner)
	// format all columns for maxlen column width
	if !ap.CmdParams.Nf {
		data.formatDataToMaxWidth(maxlen)
	}
	// print data slices as line with or without column seperator
	data.printAsciiTab()
}
