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

// T_rawdata represents raw input data as a slice of strings
type T_rawdata []string

// T_dataline represents a single line of parsed data as a slice of strings
type T_dataline []string

// T_parsedData represents the entire parsed dataset as a slice of T_dataline
type T_parsedData []T_dataline

// printJSON prints the parsed data in JSON format.
// It uses the header line defined in CmdParams.Header and the separator defined in CmdParams.Sep.
func printJSON(d T_parsedData) {
	sep := []rune(ap.CmdParams.Sep)[0]
	header := LineParse(ap.CmdParams.Header, sep)

	fmt.Println("[")
	for ln, line := range d {
		fmt.Println("  {")
		for col, val := range line {
			// Print each key-value pair
			fmt.Printf("    %q: %q", header[col], val)
			if col+1 < len(line) {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		// Close the object and add a comma if it's not the last line
		if ln+1 < len(d) {
			fmt.Println("  },")
		} else {
			fmt.Println("  }")
		}
	}
	fmt.Println("]")
}

// printJSONwithTC prints the parsed data in JSON format with a top-level collection.
// It uses the header line defined in CmdParams.Header and the separator defined in CmdParams.Sep.
func printJSONwithTC(d T_parsedData) {
	sep := []rune(ap.CmdParams.Sep)[0]
	header := d[0]
	if ap.CmdParams.Header != "" {
		header = LineParse(ap.CmdParams.Header, sep)
	}
	if ap.CmdParams.Ts {
		d = d[1:]
	}

	hkey := header[0]
	header = header[1:]
	pl := pluralize.NewClient()

	// Print the opening of the top-level collection
	fmt.Printf("{\n  %q: [\n", pl.Plural(hkey))
	for ln, line := range d {
		// Print each object in the collection
		fmt.Printf("    {\n      %q: %q,\n      \"data\": {\n", hkey, line[0])
		for col, val := range line[1:] {
			fmt.Printf("        %q: %q", header[col], val)
			if col+1 < len(line)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Println("      }")
		if ln+1 < len(d) {
			fmt.Println("    },")
		} else {
			fmt.Println("    }")
		}
	}
	// Close the top-level collection
	fmt.Println("  ]\n}")
}

// PrintJson prints the parsed data in JSON format.
// It chooses between direct JSON marshaling, printJSONwithTC, or printJSON based on CmdParams flags.
func (d T_parsedData) PrintJson() {
	// Try direct JSON marshaling if no header is specified and Ts flag is not set
	if ap.CmdParams.Header == "" && !ap.CmdParams.Ts {
		b, err := json.MarshalIndent(d, "", "  ")
		if err == nil {
			fmt.Println(string(b))
			return
		}
	}

	// Choose between printJSONwithTC and printJSON based on flags
	if ap.CmdParams.Jtc || ap.CmdParams.Ts {
		printJSONwithTC(d)
	} else {
		printJSON(d)
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

// Append appends a dataline to the parsed data
func (d *T_parsedData) Append(l T_dataline) {
	*d = append(*d, l)
}

// Insert inserts a dataline at the given position in the parsed data
func (d *T_parsedData) Insert(l T_dataline, pos int) {
	if pos >= 0 && pos <= len(*d) {
		*d = append(*d, l)
		data := *d
		copy(data[(pos+1):], data[pos:])
		data[pos] = l
		d = &data
	}
}

// generateLine formats the fields of dataline to maxlen of columns
func (data *T_dataline) generateLine(maxlen T_maxlenghts) {
	for pos, mxlen := range maxlen {
		val := ""
		if pos < len(*data) {
			val = (*data)[pos]
			runecount := utf8.RuneCountInString(val)
			blanklen := mxlen - runecount
			if regexp.MustCompile(`^ *[0-9\.,]+ *$`).MatchString(val) && !ap.CmdParams.Nn {
				(*data)[pos] = strings.Repeat(" ", blanklen) + val
			} else {
				(*data)[pos] = val + strings.Repeat(" ", blanklen)
			}
		} else {
			*data = append(*data, strings.Repeat(" ", mxlen))
		}
	}
}

// InsertGroupSeperator inserts a separator slice when the content of the specified column changes value.
func (data *T_parsedData) InsertGroupSeperator(gcol int, gcolval bool, trenner, htrenner []string) {
	if gcol <= 0 || gcol > len(*data) {
		return
	}

	gcol--
	nd := T_parsedData{}
	ref := ""

	for i, row := range *data {
		if i > 0 && len(row) > gcol && row[gcol] != trenner[gcol] && row[gcol] != htrenner[gcol] {
			if ref != row[gcol] && ref != trenner[gcol] && ref != htrenner[gcol] {
				nd.Append(trenner)
				ref = row[gcol]
				nd.Append(row)
			} else if !gcolval && ref == row[gcol] {
				row[gcol] = "''"
				nd.Append(row)
			} else {
				nd.Append(row)
			}
		} else {
			nd.Append(row)
		}
	}

	*data = nd
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
func (data *T_parsedData) insertTrenner(trenner, htrenner []string) {
	if ap.CmdParams.Ts || ap.CmdParams.Fs || ap.CmdParams.Pp {
		if ap.CmdParams.Pp {
			if ap.CmdParams.Fs {
				data.Insert(htrenner, len(*data)-1)
			}
			data.Insert(trenner, 0)
			data.Insert(trenner, 2)
			data.Append(trenner)
		} else {
			if ap.CmdParams.Ts {
				data.Insert(htrenner, 1)
			}
			if ap.CmdParams.Fs {
				data.Insert(htrenner, len(*data)-1)
			}
		}
	}
}

// sort sorts the parsed data based on the specified column index.
func (data *T_parsedData) sort(k int) {
	k--
	l1 := T_dataline{}
	d := *data

	if !ap.CmdParams.Nhl {
		l1, d = d[0], d[1:]
	}

	sort.SliceStable(d, func(i, j int) bool {
		return d[i][k] < d[j][k]
	})

	if !ap.CmdParams.Nhl {
		*data = append(T_parsedData{l1}, d...)
	} else {
		*data = d
	}
}

// delete elements from data
func (data *T_parsedData) delete(i, j int) {
	*data = slices.Delete(*data, i, j)
}

// formatDataToMaxWidth formats the data to column max width
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

// Format selects the data and header columns, inserts separators, formats the data fields,
// and prints the data as CSV, JSON, or ASCII table depending on options.
func Format(data T_parsedData) {
	sep := []rune(ap.CmdParams.Sep)[0]

	// Apply column selection if specified
	if len(ap.CmdParams.Columns) > 0 {
		data.selectColumns()
	}
	// Remove header if Rh flag is set
	if ap.CmdParams.Rh {
		data.delete(0, 1)
	}
	// Sort data if SortCol is specified
	if ap.CmdParams.SortCol > 0 {
		data.sort(int(ap.CmdParams.SortCol))
	}
	// Insert header if specified and not in JSON mode
	if ap.CmdParams.Header != "" && !ap.CmdParams.Json {
		headerline := LineParse(ap.CmdParams.Header, sep)
		if len(ap.CmdParams.Columns) > 0 && len(headerline) > len(ap.CmdParams.Columns) {
			headerline.selectColumns()
		}
		data.Insert(headerline, 0)
	}

	// Calculate maximum length for each column
	maxlen := GetMaxLength(data)
	// Insert row numbers if Num flag is set
	if ap.CmdParams.Num {
		n := make([]string, len(maxlen))
		for i := range maxlen {
			ns := strconv.Itoa(i + 1)
			if len(ap.CmdParams.Columns) > 0 {
				ns += fmt.Sprintf(" [%d]", ap.CmdParams.Columns[i])
			}
			n[i] = ns
		}
		data.Insert(n, 0)
		maxlen = GetMaxLength(data)
	}

	// Create separator lines
	trenner := make([]string, len(maxlen))
	htrenner := make([]string, len(maxlen))
	for i, v := range maxlen {
		trenner[i] = strings.Repeat("-", v)
		htrenner[i] = strings.Repeat("=", v)
	}

	// Insert separators if not in CSV or JSON mode
	if !(ap.CmdParams.Json || ap.CmdParams.Csv) {
		data.insertTrenner(trenner, htrenner)
	}

	// Output data in the appropriate format
	switch {
	case ap.CmdParams.Csv:
		data.PrintCsv()
	case ap.CmdParams.Json:
		data.PrintJson()
	default:
		data.InsertGroupSeperator(int(ap.CmdParams.Gcol), ap.CmdParams.GcolVal, trenner, htrenner)
		if !ap.CmdParams.Nf {
			data.formatDataToMaxWidth(maxlen)
		}
		data.printAsciiTab()
	}
}
