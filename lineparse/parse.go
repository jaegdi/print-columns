package pc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type T_rawdata []string
type T_dataline []string
type T_parsedData []T_dataline

func (d T_parsedData) Print() {
	b, err := json.MarshalIndent(d, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
}

func (d T_parsedData) Append(l T_dataline) T_parsedData {
	d = append(d, l)
	return d
}

func LineParse(line string, sep rune) T_dataline {
	inQuotedTextSingle := false
	inQuotedTextDouble := false
	var s string
	if sep == ' ' {
		seps := string(sep)
		re1 := regexp.MustCompile(fmt.Sprintf("%s%s+", seps, seps))
		if re1.MatchString(line) {
			re2 := regexp.MustCompile(fmt.Sprintf("([^%s])%s([^%s])", seps, seps, seps))
			s = re2.ReplaceAllString(line, `${1}_${2}`)
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

func GetLineString(line string) string {
	return fmt.Sprintln(
		LineParse(line, ' '))
}
