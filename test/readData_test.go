package main

import (
	"reflect"
	"testing"
)

// Test reading data from file
func TestReadDataFromFile(t *testing.T) {
	want := []string{
		"A B C D E F",
		"F  B  C  A B  E  F",
		"a bb ccc dd eeeee ffffff",
	}

	erg, err := readTestDataFile("data.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`readTestDataFile("data.txt") = %q, want %q`, erg, want)
	}
}
