package pc

import (
	"fmt"
	ap "pc/argparse"
	"regexp"
	"strings"
)

// LineParse parses the text line and split into fields, return a slice of strings.
// optionale it can detect cloiuns with blanks, but therefor the columns in the textline
// must separated by more than one blank.
func LineParse(line string, sep rune) T_dataline {
	inQuotedTextSingle := false
	inQuotedTextDouble := false
	var s string
	if sep == ' ' && ap.CmdParams.MoreBlanks {
		seps := string(sep)
		re1 := regexp.MustCompile(fmt.Sprintf("%s%s+", seps, seps))
		if re1.MatchString(line) {
			re2 := regexp.MustCompile(fmt.Sprintf("([^%s])%s([^%s])", seps, seps, seps))
			s = re2.ReplaceAllString(line, `${1}§${2}`)
			s = re1.ReplaceAllString(s, seps)
		} else {
			s = line
		}
	} else {
		s = line
	}
	ss := T_dataline(strings.FieldsFunc(s, func(r rune) bool {
		if r == sep && !inQuotedTextSingle && !inQuotedTextDouble {
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
	return ss
}

// GetLineSlice returns a parsed and splittet line as string
func GetLineSlice(line string) string {
	return fmt.Sprintln(
		LineParse(line, ' '))
}

// DataParse parses an slice of stringlines into T_parsedData ( [][]string )
func DataParse(data T_rawdata, sep rune) T_parsedData {
	pdata := T_parsedData{}
	for _, l := range data {
		dataline := LineParse(l, sep)
		// fmt.Println(dataline)
		pdata = append(pdata, dataline)
	}
	return pdata
}
