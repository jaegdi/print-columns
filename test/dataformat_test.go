package main

import (
	"strings"
	"testing"

	ap "pc/internal/argparse"
	df "pc/internal/dataformat"
)

func TestPrintJSONPart1(t *testing.T) {
	data := df.T_parsedData{
		df.T_dataline{"John", "Doe", "john.doe@example.com"},
		df.T_dataline{"Jane", "Smith", "jane.smith@example.com"},
	}
	// Test case 1: Print JSON with default header and separator
	expectedOutput1txt := `[
  {
    "Prename": "John",
    "Surname": "Doe",
    "Email": "john.doe@example.com"
  },
  {
    "Prename": "Jane",
    "Surname": "Smith",
    "Email": "jane.smith@example.com"
  }
]`
	ap.CmdParams.Header = "Prename Surname Email"
	ap.CmdParams.Sep = " "
	actualOutput1 := df.PrintJSON(data)
	// if reflect.DeepEqual(expectedOutput1txt, actualOutput1) {
	if strings.Compare(strings.TrimSpace(expectedOutput1txt), strings.TrimSpace(actualOutput1)) != 0 {
		t.Errorf("PrintJSON output does not match expected.\nExpected:\n%s\nActual:\n%s", strings.TrimSpace(expectedOutput1txt), strings.TrimSpace(actualOutput1))
	}
}

func TestPrintJSONPart2(t *testing.T) {
	data := df.T_parsedData{
		df.T_dataline{"John", "Doe", "john.doe@example.com"},
		df.T_dataline{"Jane", "Smith", "jane.smith@example.com"},
	}

	// Test case 2: Print JSON with custom header and separator
	ap.CmdParams.Header = "First Name,Last Name,Email"
	ap.CmdParams.Sep = ","
	expectedOutput2txt := `[
  {
    "First Name": "John",
    "Last Name": "Doe",
    "Email": "john.doe@example.com"
  },
  {
    "First Name": "Jane",
    "Last Name": "Smith",
    "Email": "jane.smith@example.com"
  }
]`
	actualOutput2 := df.PrintJSON(data)
	if strings.Compare(strings.TrimSpace(expectedOutput2txt), strings.TrimSpace(actualOutput2)) != 0 {
		t.Errorf("PrintJSON output does not match expected.\nExpected:\n%s\nActual:\n%s", strings.TrimSpace(expectedOutput2txt), strings.TrimSpace(actualOutput2))
	}
}
