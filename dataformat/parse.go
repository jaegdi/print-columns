package pc

import (
	"fmt"
	ap "pc/argparse"
	"regexp"
	"strconv"
	"strings"
)

// LineParse parses the text line and split into fields, return a slice of strings.
// optionale it can detect columns with blanks, but therefor the columns in the textline
// must separated by more than one blank.
func LineParse(line string, sep rune) T_dataline {
	inQuotedTextSingle := false
	inQuotedTextDouble := false
	var s string
	// if MoreBlanks is set, then only two sep are recognized as columns sep, a single sep
	// is replaced by '\n' (this can not be content of a line)
	if sep == ' ' && ap.CmdParams.MoreBlanks {
		seps := string(sep)
		re1 := regexp.MustCompile(fmt.Sprintf("%s%s+", seps, seps))
		if re1.MatchString(line) {
			re2 := regexp.MustCompile(fmt.Sprintf("([^%s])%s([^%s])", seps, seps, seps))
			s = re2.ReplaceAllString(line, "${1}\n${2}")
			s = re1.ReplaceAllString(s, seps)
		} else {
			s = line
		}
	} else {
		s = line
	}
	ss := T_dataline(strings.FieldsFunc(s, func(r rune) bool {
		//  When sep = ' ' \t is also accepted as sep
		if (r == sep || (sep == ' ' && r == '\t')) && !inQuotedTextSingle && !inQuotedTextDouble {
			return true
		}
		if r == '\'' {
			inQuotedTextSingle = !inQuotedTextSingle
		}
		if r == '"' {
			inQuotedTextDouble = !inQuotedTextDouble
		}
		return false
	}))
	for i := range ss {
		ss[i] = strings.Replace(ss[i], "\n", " ", -1)
	}
	return ss
}

// GetLineSlice returns a parsed and splittet line as string
func GetLineSlice(line string) string {
	return fmt.Sprintln(
		LineParse(line, ' '))
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
