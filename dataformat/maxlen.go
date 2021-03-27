package pc

import (
	"unicode/utf8"
)

type T_maxlenghts []int

func (ml *T_maxlenghts) setMl(col int, length int) {
	if len(*ml) < col+1 {
		*ml = append(*ml, 0)
	}
	m := *ml
	// fmt.Println(ml)
	if m[col] < length {
		m[col] = length
	}
	ml = &m
}

func GetMaxLength(d T_parsedData) T_maxlenghts {
	maxlengths := T_maxlenghts{}
	for _, line := range d {
		for col, val := range line {
			length := utf8.RuneCountInString(val)
			maxlengths.setMl(col, length)
		}
	}
	return maxlengths
}
