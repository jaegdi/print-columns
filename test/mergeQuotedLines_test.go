package main

import (
	"reflect"
	"testing"

	ld "pc/loaddata"
)

func TestMergeQuotedLines(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "no quotes",
			input: []string{
				"line1 col1 col2",
				"line2 col1 col2",
			},
			expected: []string{
				"line1 col1 col2",
				"line2 col1 col2",
			},
		},
		{
			name: "single line with quotes",
			input: []string{
				`col1 "quoted value" col3`,
			},
			expected: []string{
				`col1 "quoted value" col3`,
			},
		},
		{
			name: "multiline quoted field",
			input: []string{
				`col1 "this is a`,
				`multiline value" col3`,
			},
			expected: []string{
				"col1 \"this is a\nmultiline value\" col3",
			},
		},
		{
			name: "three line quoted field",
			input: []string{
				`AlertName openshift-logging "The message`,
				`continues here`,
				`and ends here"`,
			},
			expected: []string{
				"AlertName openshift-logging \"The message\ncontinues here\nand ends here\"",
			},
		},
		{
			name: "mixed single and multiline",
			input: []string{
				`single1 "value1"`,
				`multi "starts here`,
				`ends here"`,
				`single2 "value2"`,
			},
			expected: []string{
				`single1 "value1"`,
				"multi \"starts here\nends here\"",
				`single2 "value2"`,
			},
		},
		{
			name: "multiple quotes on same line",
			input: []string{
				`col1 "val1" "val2" col4`,
			},
			expected: []string{
				`col1 "val1" "val2" col4`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ld.MergeQuotedLines(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeQuotedLines() = %v, want %v", result, tt.expected)
			}
		})
	}
}
