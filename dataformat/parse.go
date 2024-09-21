package pc

import (
	"fmt"
	ap "pc/argparse"
	"regexp"
	"strconv"
	"strings"
)

// LineParse splits a text line into fields, optionally handling columns with multiple spaces.
// It returns a slice of strings (T_dataline).
func LineParse(line string, sep rune) T_dataline {
	// Handle the case of multiple spaces as separators
	if sep == ' ' && ap.CmdParams.MoreBlanks {
		line = handleMultipleSpaces(line)
	}

	// Split the line into fields
	fields := splitFields(line, sep)

	// Replace any newline characters with spaces in each field
	for i := range fields {
		fields[i] = strings.ReplaceAll(fields[i], "\n", " ")
	}

	return T_dataline(fields)
}

// handleMultipleSpaces replaces single spaces with newlines and multiple spaces with a single space.
func handleMultipleSpaces(s string) string {
	re1 := regexp.MustCompile(`\s{2,}`)
	re2 := regexp.MustCompile(`([^\s])\s([^\s])`)

	s = re1.ReplaceAllString(s, " ")
	return re2.ReplaceAllString(s, "${1}\n${2}")
}

// splitFields splits a string into fields based on the separator, respecting quoted text.
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

// GetLineSlice returns a parsed and splittet line as string
func GetLineSlice(line string) string {
	return fmt.Sprintln(LineParse(line, ' '))
}

// setFilter returns the col for filter (-1 if col undefined) and the compiled regexp
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

// DataParse parses an slice of stringlines into T_parsedData ( [][]string )
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
