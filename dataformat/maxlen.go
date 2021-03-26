package pc

import (
	lp "pc/lineparse"
)

type T_maxlenghts []int

func setMl(d T_maxlenghts, col int, length int) T_maxlenghts {
	if len(d) < col+1 {
		d = append(d, 0)
	}
	if d[col] < length {
		d[col] = length
	}
	return d
}

func GetMaxLength(d lp.T_parsedData) T_maxlenghts {
	maxlengths := T_maxlenghts{}
	for _, line := range d {
		for col, val := range line {
			length := len(val)
			setMl(maxlengths, col, length)
		}
	}
	return maxlengths
}
