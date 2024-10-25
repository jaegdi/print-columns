package pc

import (
	"unicode/utf8"
)

// T_maxlenghts represents a slice of integers storing maximum lengths.
type T_maxlenghts []int

/**
 * @brief Sets the maximum length for a given column index.
 *
 * Updates the maximum length stored in the T_maxlenghts slice for a specific column index. If the slice is too short, it's extended.
 * @param col The column index (starting from 0).
 * @param length The length to set.
 */
// setMl sets the maximum length for a given column index.
// It takes a column index `col` and a length `length` as input, and updates the T_maxlenghts slice accordingly.  If the slice is too short, it's extended.
func (ml *T_maxlenghts) setMl(col int, length int) {
	if len(*ml) < col+1 {
		*ml = append(*ml, 0)
	}
	if (*ml)[col] < length {
		(*ml)[col] = length
	}
}

/**
 * @brief Calculates the maximum length of each column in the parsed data.
 *
 * Iterates through the parsed data to determine the maximum length of each column.  Returns a slice of integers representing the maximum length of each column.
 * @param d The parsed data.
 * @return T_maxlenghts A slice of integers representing the maximum length of each column.
 */
// GetMaxLength calculates the maximum length of each column in the parsed data.
// It iterates through each row and column, updating the maximum length for each column as needed.
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
