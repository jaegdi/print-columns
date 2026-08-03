package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	ap "pc/argparse"
	df "pc/dataformat"
)

// captureOutput captures stdout during the execution of a function
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// resetCmdParams resets the command parameters to defaults for testing
func resetCmdParams() {
	ap.CmdParams.Pp = false
	ap.CmdParams.Cs = false
	ap.CmdParams.Colsep = "|"
	ap.CmdParams.ColSepW = 1
	ap.CmdParams.MoreBlanks = false
	ap.CmdParams.Mark = ""
	ap.CmdParams.Nf = false
	ap.CmdParams.Sep = " "
	ap.CmdParams.Header = ""
	ap.CmdParams.Columns = nil
	ap.CmdParams.Rh = false
	ap.CmdParams.SortCol = 0
	ap.CmdParams.Num = false
	ap.CmdParams.Json = false
	ap.CmdParams.Csv = false
	ap.CmdParams.Ts = false
	ap.CmdParams.Fs = false
	ap.CmdParams.Gcol = 0
	ap.CmdParams.Nhl = false
}

func TestPrintAsciiTabMultilineCell(t *testing.T) {
	resetCmdParams()
	ap.CmdParams.Pp = true

	// Create test data with a multiline cell (Description contains a linefeed)
	data := df.T_parsedData{
		df.T_dataline{"Name", "Description", "Value"},
		df.T_dataline{"Alice", "Hello\nWorld", "100"},
		df.T_dataline{"Bob", "Simple", "200"},
	}

	output := captureOutput(func() {
		df.Format(data)
	})

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// With -pp enabled, we expect:
	// Line 0: separator  | ----- | ----------- | ----- |
	// Line 1: header     | Name  | Description | Value |
	// Line 2: separator  | ----- | ----------- | ----- |
	// Line 3: Alice row line 1: | Alice | Hello       |   100 |
	// Line 4: Alice row line 2: |       | World       |       |
	// Line 5: Bob row:   | Bob   | Simple      |   200 |
	// Line 6: separator  | ----- | ----------- | ----- |

	// Find the Alice row - it should span 2 visual lines
	aliceLineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Alice") {
			aliceLineIdx = i
			break
		}
	}

	if aliceLineIdx < 0 {
		t.Fatalf("Could not find Alice row in output:\n%s", output)
	}

	// Check that Alice row line 1 contains "Hello" and "100"
	if !strings.Contains(lines[aliceLineIdx], "Hello") {
		t.Errorf("Alice row line 1 should contain 'Hello', got: %s", lines[aliceLineIdx])
	}
	if !strings.Contains(lines[aliceLineIdx], "100") {
		t.Errorf("Alice row line 1 should contain '100', got: %s", lines[aliceLineIdx])
	}

	// Check that the next line (Alice row line 2) contains "World" and has blank padding for other cells
	aliceLine2 := lines[aliceLineIdx+1]
	if !strings.Contains(aliceLine2, "World") {
		t.Errorf("Alice row line 2 should contain 'World', got: %s", aliceLine2)
	}
	// Should NOT contain "Alice" or "100" on the continuation line
	if strings.Contains(aliceLine2, "Alice") {
		t.Errorf("Alice row line 2 should not repeat 'Alice', got: %s", aliceLine2)
	}

	// Check that Bob row follows after Alice's multiline
	bobLineIdx := aliceLineIdx + 2
	if bobLineIdx >= len(lines) || !strings.Contains(lines[bobLineIdx], "Bob") {
		t.Errorf("Bob row should follow Alice's multiline row. Output:\n%s", output)
	}

	t.Logf("Output:\n%s", output)
}

func TestPrintAsciiTabMultipleCellsWithLinefeeds(t *testing.T) {
	resetCmdParams()
	ap.CmdParams.Pp = true

	// Create test data where multiple cells have linefeeds
	data := df.T_parsedData{
		df.T_dataline{"Name", "Desc", "Notes"},
		df.T_dataline{"Alice", "Line1\nLine2\nLine3", "A\nB"},
	}

	output := captureOutput(func() {
		df.Format(data)
	})

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// Find Alice row
	aliceLineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Alice") {
			aliceLineIdx = i
			break
		}
	}

	if aliceLineIdx < 0 {
		t.Fatalf("Could not find Alice row in output:\n%s", output)
	}

	// Alice's row should span 3 visual lines (max of 3 in Desc, 2 in Notes)
	// Line 1: Alice | Line1 | A
	// Line 2:       | Line2 | B
	// Line 3:       | Line3 |   (blank padding)

	if !strings.Contains(lines[aliceLineIdx], "Alice") || !strings.Contains(lines[aliceLineIdx], "Line1") || !strings.Contains(lines[aliceLineIdx], "A") {
		t.Errorf("Alice row line 1 incorrect: %s", lines[aliceLineIdx])
	}

	if aliceLineIdx+1 >= len(lines) || !strings.Contains(lines[aliceLineIdx+1], "Line2") || !strings.Contains(lines[aliceLineIdx+1], "B") {
		t.Errorf("Alice row line 2 incorrect: %s", lines[aliceLineIdx+1])
	}

	if aliceLineIdx+2 >= len(lines) || !strings.Contains(lines[aliceLineIdx+2], "Line3") {
		t.Errorf("Alice row line 3 should contain 'Line3': %s", lines[aliceLineIdx+2])
	}

	t.Logf("Output:\n%s", output)
}

func TestPrintAsciiTabNoLinefeeds(t *testing.T) {
	resetCmdParams()
	ap.CmdParams.Pp = true

	// Standard data without linefeeds - should work as before
	data := df.T_parsedData{
		df.T_dataline{"Name", "Value"},
		df.T_dataline{"Alice", "100"},
		df.T_dataline{"Bob", "200"},
	}

	output := captureOutput(func() {
		df.Format(data)
	})

	// Should have exactly 6 lines (top border, header, separator, Alice, Bob, bottom border)
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != 6 {
		t.Errorf("Expected 6 lines for standard table, got %d:\n%s", len(lines), output)
	}

	t.Logf("Output:\n%s", output)
}

func TestPrintAsciiTabQuotedMultiline(t *testing.T) {
	resetCmdParams()
	ap.CmdParams.Pp = true

	// Test with quoted value containing linefeed: '"line1\nline2"'
	data := df.T_parsedData{
		df.T_dataline{"Name", "Description"},
		df.T_dataline{"Alice", "\"line1\nline2\""},
		df.T_dataline{"Bob", "simple"},
	}

	output := captureOutput(func() {
		df.Format(data)
	})

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// Find Alice row
	aliceLineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Alice") {
			aliceLineIdx = i
			break
		}
	}

	if aliceLineIdx < 0 {
		t.Fatalf("Could not find Alice row in output:\n%s", output)
	}

	// Alice row should span 2 visual lines
	// Line 1: Alice | "line1
	// Line 2:       | line2"
	if !strings.Contains(lines[aliceLineIdx], "Alice") || !strings.Contains(lines[aliceLineIdx], "\"line1") {
		t.Errorf("Alice row line 1 should contain 'Alice' and '\"line1', got: %s", lines[aliceLineIdx])
	}

	if aliceLineIdx+1 >= len(lines) || !strings.Contains(lines[aliceLineIdx+1], "line2\"") {
		t.Errorf("Alice row line 2 should contain 'line2\"', got: %s", lines[aliceLineIdx+1])
	}

	// Bob should follow after Alice's multiline
	bobLineIdx := aliceLineIdx + 2
	if bobLineIdx >= len(lines) || !strings.Contains(lines[bobLineIdx], "Bob") {
		t.Errorf("Bob row should follow Alice's multiline row. Output:\n%s", output)
	}

	t.Logf("Output:\n%s", output)
}

func TestPrintAsciiTabWithoutPrettyPrint(t *testing.T) {
	resetCmdParams()
	// No -pp flag, just plain output

	data := df.T_parsedData{
		df.T_dataline{"Name", "Description"},
		df.T_dataline{"Alice", "Hello\nWorld"},
	}

	output := captureOutput(func() {
		df.Format(data)
	})

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")

	// Find Alice row
	aliceLineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Alice") {
			aliceLineIdx = i
			break
		}
	}

	if aliceLineIdx < 0 {
		t.Fatalf("Could not find Alice row in output:\n%s", output)
	}

	// Should still handle multiline
	if !strings.Contains(lines[aliceLineIdx], "Hello") {
		t.Errorf("Alice row line 1 should contain 'Hello': %s", lines[aliceLineIdx])
	}
	if aliceLineIdx+1 >= len(lines) || !strings.Contains(lines[aliceLineIdx+1], "World") {
		t.Errorf("Alice row line 2 should contain 'World': %s", output)
	}

	t.Logf("Output:\n%s", output)
}
