package pc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	ap "pc/internal/argparse"
)

/**
 * @brief Splits a text line into fields, optionally handling columns with multiple spaces.
 *
 * Parses a single line of text into a slice of strings (T_dataline), using either a specified separator or handling multiple spaces as separators if the `MoreBlanks` flag is set.
 * @param line The input line of text.
 * @param sep The separator character to use (if `MoreBlanks` is false).
 * @return T_dataline A slice of strings representing the parsed fields.
 */
// LineParse splits a text line into fields, optionally handling columns with multiple spaces.
// It returns a slice of strings (T_dataline).
func LineParse(line string, sep rune) T_dataline {

	var fields []string

	// Handle the case of multiple spaces as separators
	if sep == ' ' && ap.CmdParams.MoreBlanks {
		line = handleMultipleSpaces(line)
		fields = splitFields(line, '\n')
	} else {
		fields = splitFields(line, sep)
	}

	return T_dataline(fields)
}

/**
 * @brief Replaces multiple spaces with a single newline character.
 *
 * Replaces consecutive spaces in a string with a single newline character.  Used for parsing lines with multiple spaces as delimiters.
 * @param s The input string.
 * @return string The modified string with single newlines replacing multiple spaces.
 */
// handleMultipleSpaces replaces multiple spaces with a single space and ensures single-character fields are recognized.
// It takes a string `s` as input and returns a modified string.
func handleMultipleSpaces(s string) string {
	// Replace multiple spaces with a single space
	re1 := regexp.MustCompile(`\s{2,}`)
	s = re1.ReplaceAllString(s, "\n")

	return s
}

/**
 * @brief Splits a string into fields based on a separator, handling quoted text.
 *
 * Splits a string into a slice of strings based on a given separator character, correctly handling quoted text to prevent splitting within quotes.
 * @param s The input string.
 * @param sep The separator character.
 * @return []string A slice of strings representing the split fields.
 */
// splitFields splits a string into fields based on the separator, respecting quoted text.
// It takes a string `s` and a rune `sep` as input and returns a slice of strings.
func splitFields(s string, sep rune) []string {
	var fields []string
	var field strings.Builder
	inQuoteSingle, inQuoteDouble := false, false

	for _, r := range s {
		switch {
		case r == '\'' && !inQuoteDouble:
			// Toggle single quote state if not in double quotes
			inQuoteSingle = !inQuoteSingle
			field.WriteRune(r)
		case r == '"' && !inQuoteSingle:
			// Toggle double quote state if not in single quotes
			inQuoteDouble = !inQuoteDouble
			field.WriteRune(r)
		case (r == sep || (sep == ' ' && r == '\t')) && !inQuoteSingle && !inQuoteDouble:
			// If separator is found and not in quotes, add field to fields and reset
			if field.Len() > 0 {
				fields = append(fields, field.String())
				field.Reset()
			}
		default:
			// Add character to current field
			field.WriteRune(r)
		}
	}

	if field.Len() > 0 {
		fields = append(fields, field.String())
	}

	return fields
}

/**
 * @brief Returns a parsed and split line as a string.
 * @param line The input line.
 * @return string The formatted string.
 */
// GetLineSlice returns a parsed and split line as a string.
// It takes a string `line` as input and returns a formatted string.
func GetLineSlice(line string) string {
	return fmt.Sprintln(LineParse(line, ' '))
}

/**
 * @brief Sets the filter parameters for data parsing.
 *
 * Parses the filter string from command-line parameters to determine the column index and regular expression pattern for filtering data.
 * @return (int, *regexp.Regexp) The column index for filtering (-1 if no column is specified) and the compiled regular expression.
 */
// setFilter returns the column index for filtering (-1 if no column is specified) and the compiled regular expression.
// It parses the filter string from command-line parameters to determine the column index and regular expression pattern.
func setFilter() (int, *regexp.Regexp) {
	col := -1
	filterString := ap.CmdParams.Filter
	if ap.CmdParams.Filter != "" {
		filtertRegExp := regexp.MustCompile(`^(\d+)=(.*)$`)
		if filtertRegExp.MatchString(ap.CmdParams.Filter) {
			res := filtertRegExp.FindAllStringSubmatch(ap.CmdParams.Filter, -1)
			i, err := strconv.Atoi(res[0][1])
			if err == nil {
				col = i - 1
				filterString = res[0][2]
				fmt.Println("col:", col, "pattern", filterString)
			}
		}
	} else {
		filterString = `.`
	}
	dataRegExp := regexp.MustCompile(filterString)
	return col, dataRegExp
}

/**
 * @brief Parses a slice of string lines into a T_parsedData structure, applying filtering if specified.
 *
 * Parses a slice of raw data strings into a structured T_parsedData, applying filtering based on a provided regular expression and column index if a filter string is specified in the command-line parameters.
 * @param data The raw data as a slice of strings.
 * @param sep The separator character.
 * @return T_parsedData The parsed data.
 */
// DataParse parses a slice of string lines into a T_parsedData structure.
// It applies filtering based on command-line parameters if specified.
func DataParse(data T_rawdata, sep rune) T_parsedData {
	pdata := T_parsedData{}
	var filterRegExp *regexp.Regexp
	filterCol := -1
	filter := ap.CmdParams.Filter != ""
	if filter {
		filterCol, filterRegExp = setFilter()
	}
	for _, l := range data {
		var dataline T_dataline
		if !filter || (filterCol < 0 && filterRegExp.MatchString(l)) {
			dataline = LineParse(l, sep)
			pdata = append(pdata, dataline)
		} else {
			dataline = LineParse(l, sep)
			if filterCol > -1 && filterCol < len(dataline) && filterRegExp.MatchString(dataline[filterCol]) {
				pdata = append(pdata, dataline)
			}
		}
	}
	return pdata
}
