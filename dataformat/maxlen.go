package pc

import (
	"unicode/utf8"
)

type T_maxlenghts []int

// setMl set the max len in T_maxlengths
func (ml *T_maxlenghts) setMl(col int, length int) {
	if len(*ml) < col+1 {
		*ml = append(*ml, 0)
	}
	// m := *ml
	// fmt.Println(ml)
	if (*ml)[col] < length {
		(*ml)[col] = length
	}
	// ml = &m
}

//  GetMaxLength check the columns of all rows and get the max length for each column
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
