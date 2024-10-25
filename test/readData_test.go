package main

import (
	"reflect"
	"testing"

	ld "pc/internal/loaddata"
)

// Test GetLineString, line with fixed length columns seperated by two or more blanks
func TestReadDataFromFile(t *testing.T) {
	filename := "./data/data.txt"
	want := []string{
		"A  B   9  C    D   E      F       3",
		"F  B   5  C    A   B      E       F",
		"a  bb  2  ccc  dd  eeeee  ffffff  x",
		"1  2   1  3    2   5      6       1",
	}
	erg, err := ld.GetData(filename)
	if err != nil {
		t.Fatalf("GetData(\"%s\") returned error: %v", filename, err)
	}
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q`, filename, erg, want)
	}
}

func TestReadEmptyFile(t *testing.T) {
	filename := "./data/empty.txt"
	want := []string{}
	erg, err := ld.GetData(filename)
	if err != nil {
		t.Fatalf("GetData(\"%s\") returned error: %v", filename, err)
	}
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q`, filename, erg, want)
	}

}

func TestReadNonExistentFile(t *testing.T) {
	filename := "./data/nonexistent.txt"
	_, err := ld.GetData(filename)
	if err == nil {
		t.Fatalf("GetData(%q) did not return an error for a nonexistent file", filename)
	}
}

func TestReadDifferentLineEndings(t *testing.T) {
	//This test requires creating a file with different line endings.  This is outside the scope of this task.
}
