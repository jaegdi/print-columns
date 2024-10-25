package pc

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	pluralize "github.com/gertd/go-pluralize"
	"golang.org/x/exp/slices"

	ap "pc/internal/argparse"
)

// T_rawdata represents raw input data as a slice of strings.
type T_rawdata []string

// T_dataline represents a single line of parsed data as a slice of strings.
type T_dataline []string

// T_parsedData represents the entire parsed dataset as a slice of T_dataline.
type T_parsedData []T_dataline

/**
 * @brief Prints a horizontal separator line to the console.
 */
func horLine() {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()
}

/**
 * @brief Prints data to the console with a message.
 * @param msg The message to print before the data.
 * @param data The data to print.  Supports various data types including slices of strings and custom data structures.
 */
func PrintData(msg string, data interface{}) {
	fmt.Println(msg)
	switch d := data.(type) {
	case []string:
		for i, line := range d {
			fmt.Println(i, "[]string)", ":", line)
		}
	case []T_dataline:
		for i, line := range d {
			fmt.Println(i, "[]T_dataline", ":", line)
		}
	case T_parsedData:
		for i, line := range d {
			fmt.Println(i, "T_parsedData", ":", line)
		}
	case T_maxlenghts:
		for i, line := range d {
			fmt.Println(i, "T_maxlenghts", ":", line)
		}
	default:
		fmt.Println(data)
	}
	horLine()
}

/**
 * @brief Converts data to a string with a message.
 * @param msg The message to prepend to the string representation of the data.
 * @param data The data to convert to a string. Supports various data types including slices of strings and custom data structures.
 * @return string The string representation of the data.
 */
func SprintData(msg string, data interface{}) string {
	s := msg + "\n"
	switch d := data.(type) {
	case []string:
		for i, line := range d {
			s += fmt.Sprintln(i, "[]string)", ":", line)
		}
	case []T_dataline:
		for i, line := range d {
			s += fmt.Sprintln(i, "[]T_dataline", ":", line)
		}
	case T_dataline:
		s += fmt.Sprintln("T_dataline", ":", data)
	case T_parsedData:
		for i, line := range d {
			s += fmt.Sprintln(i, "T_parsedData", ":", line)
		}
	case T_maxlenghts:
		for i, line := range d {
			s += fmt.Sprintln(i, "T_maxlenghts", ":", line)
		}
	default:
		s += fmt.Sprintln(data)
	}
	s += `---------------------------------------------`
	return s
}

/**
 * @brief Prints the parsed data in JSON format.
 *
 * Formats the parsed data as a JSON array of objects.  Each object represents a row, with keys derived from the header line (if provided).
 * @param d The parsed data to print.
 * @return string The JSON formatted string.
 */
// PrintJSON prints the parsed data in JSON format.
// It uses the header line defined in CmdParams.Header and the separator defined in CmdParams.Sep.
func PrintJSON(d T_parsedData) string {
	// log.Println("param d:", d)
	sep := []rune(ap.CmdParams.Sep)[0]
	// log.Println("Header: ", ap.CmdParams.Header)
	header := LineParse(ap.CmdParams.Header, sep)
	output := ""
	fmt.Println("[")
	output += fmt.Sprintln("[")
	for ln, line := range d {
		fmt.Println("  {")
		output += fmt.Sprintln("  {")
		for col, val := range line {
			// Print each key-value pair
			if header != nil {
				fmt.Printf("    %q: %q", header[col], val)
				output += fmt.Sprintf("    %q: %q", header[col], val)
			} else {
				fmt.Printf("    %q: %q", strconv.Itoa(col), val)
				output += fmt.Sprintf("    %q: %q", strconv.Itoa(col), val)
			}
			if col+1 < len(line) {
				fmt.Println(",")
				output += fmt.Sprintln(",")
			} else {
				fmt.Println()
				output += fmt.Sprintln()
			}
		}
		// Close the object and add a comma if it's not the last line
		if ln+1 < len(d) {
			fmt.Println("  },")
			output += fmt.Sprintln("  },")
		} else {
			fmt.Println("  }")
			output += fmt.Sprintln("  }")
		}
	}
	fmt.Println("]")
	output += fmt.Sprintln("]")
	return output
}

/**
 * @brief Prints the parsed data in JSON format with a top-level collection.
 *
 * Formats the parsed data as a JSON object with a top-level array. The first column is used as the key for each object in the array.
 * @param d The parsed data to print.
 */
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

/**
 * @brief Prints the parsed data in JSON format.
 *
 * Selects the appropriate JSON printing function based on command-line flags.  Handles cases with and without headers, and with and without the top-level collection.
 */
// PrintJson prints the parsed data in JSON format.
// It chooses between direct JSON marshaling, printJSONwithTC, or PrintJSON based on CmdParams flags.
func (d T_parsedData) PrintJson() {
	// Try direct JSON marshaling if no header is specified and Ts flag is not set
	if ap.CmdParams.Header == "" && !ap.CmdParams.Ts {
		b, err := json.MarshalIndent(d, "", "  ")
		if err == nil {
			fmt.Println(string(b))
			return
		}
	}

	// Choose between printJSONwithTC and PrintJSON based on flags
	if ap.CmdParams.Jtc || ap.CmdParams.Ts {
		printJSONwithTC(d)
	} else {
		_ = PrintJSON(d)
	}
}

/**
 * @brief Prints the parsed data in CSV format.
 *
 * Writes the parsed data to standard output in CSV format using the encoding/csv package.
 * @param d The parsed data to print.
 */
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

/**
 * @brief Appends a data line to the parsed data.
 * @param l The data line to append.
 */
// Append appends a dataline to the parsed data.
func (d *T_parsedData) Append(l T_dataline) {
	*d = append(*d, l)
}

/**
 * @brief Inserts a data line at a specified position in the parsed data.
 * @param l The data line to insert.
 * @param pos The position at which to insert the data line.
 */
// Insert inserts a dataline at the given position in the parsed data.
func (d *T_parsedData) Insert(l T_dataline, pos int) {
	*d = slices.Insert(*d, pos, l)
}

/**
 * @brief Formats the fields of a data line to match the maximum column lengths.
 * @param maxlen A slice of integers representing the maximum length of each column.
 */
// generateLine formats the fields of dataline to maxlen of columns.
func (data *T_dataline) generateLine(maxlen T_maxlenghts) {
	for pos, mxlen := range maxlen {
		val := ""
		if pos < len(*data) {
			val = (*data)[pos]
			runecount := utf8.RuneCountInString(val)
			blanklen := mxlen - runecount
			if regexp.MustCompile(`^ *[0-9\.,]+(k|m|d|h|H|M|J|Y|Ki|Mi|Gi){0,1} *$`).MatchString(val) && !ap.CmdParams.Nn {
				(*data)[pos] = strings.Repeat(" ", blanklen) + val
			} else {
				(*data)[pos] = val + strings.Repeat(" ", blanklen)
			}
		} else {
			*data = append(*data, strings.Repeat(" ", mxlen))
		}
	}
}

/**
 * @brief Inserts group separators into the parsed data based on changes in a specified column.
 * @param gcol The index of the column to check for changes.
 * @param gcolval A boolean indicating whether to replace values in the group column.
 * @param trenner A slice of strings representing the separator line.
 * @param htrenner A slice of strings representing the header separator line.
 */
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

/**
 * @brief Selects data columns from the parsed data based on CmdParams.Columns.
 */
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

/**
 * @brief Selects columns from a data line based on CmdParams.Columns.
 */
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

/**
 * @brief Inserts separators for TitleSeparator, FooterSeparator, or PrettyPrint.
 * @param trenner Separator line.
 * @param htrenner Header separator line.
 */
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

/**
 * @brief Checks if a slice of strings contains only numeric values.
 * @param s The slice of strings to check.
 * @return bool True if all strings are numeric, false otherwise.
 */
// isNumeric checks if a slice of strings contains only numeric values.
func isNumeric(s []string) bool {
	for _, str := range s {
		_, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return false
		}
	}
	return true
}

/**
 * @brief Sorts the parsed data based on the specified column index.
 * @param k The index of the column to sort by.
 */
// sort sorts the parsed data based on the specified column index.
func (data *T_parsedData) sort(k int) {
	k--
	l1 := T_dataline{}
	d := *data

	if !ap.CmdParams.Nhl {
		l1, d = d[0], d[1:]
	}

	// Check if the column contains only numeric values
	allNumeric := isNumeric(getColumn(d, k))

	sort.SliceStable(d, func(i, j int) bool {
		if allNumeric {
			// Numeric sort
			numI, _ := strconv.ParseFloat(d[i][k], 64)
			numJ, _ := strconv.ParseFloat(d[j][k], 64)
			return numI < numJ
		} else {
			// Lexicographical sort
			return d[i][k] < d[j][k]
		}
	})

	if !ap.CmdParams.Nhl {
		*data = append(T_parsedData{l1}, d...)
	} else {
		*data = d
	}
}

/**
 * @brief Extracts a specific column from a 2D string slice.
 * @param data The 2D string slice.
 * @param colIndex The index of the column to extract.
 * @return []string The extracted column as a slice of strings.
 */
// getColumn extracts a specific column from a 2D string slice.
func getColumn(data T_parsedData, colIndex int) []string {
	column := make([]string, len(data))
	for i, row := range data {
		if colIndex < len(row) {
			column[i] = row[colIndex]
		}
	}
	return column
}

/**
 * @brief Deletes elements from the parsed data.
 * @param i The starting index of the elements to delete.
 * @param j The ending index of the elements to delete.
 */
// delete elements from data.
func (data *T_parsedData) delete(i, j int) {
	*data = slices.Delete(*data, i, j)
}

/**
 * @brief Formats the data to the maximum column width.
 * @param maxlen A slice of integers representing the maximum length of each column.
 */
// formatDataToMaxWidth formats the data to column max width.
func (data *T_parsedData) formatDataToMaxWidth(maxlen []int) {
	for i := range *data {
		(*data)[i].generateLine(maxlen)
	}
}

/**
 * @brief Prints the parsed data as an ASCII table.
 *
 * Formats and prints the parsed data to the console as an ASCII table, applying formatting options such as column separators, widths, and borders based on command-line flags.
 * @param data The parsed data to print.
 */
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

/**
 * @brief Formats and prints the parsed data.
 *
 * This function orchestrates the data formatting and output process, selecting the appropriate output format (CSV, JSON, or ASCII table) and applying various formatting options based on command-line flags.
 * @param data The parsed data to format and print.
 */
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
