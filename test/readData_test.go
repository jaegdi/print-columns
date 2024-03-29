package main

import (
	ap "pc/argparse"
	ld "pc/loaddata"
	"reflect"
	"testing"
)

// Test GetLineString, line with fixed length columns seperated by two or more blanks
func TestReadDataFromFile(t *testing.T) {
	filename := `/home/jaegdi/devel/go/pc-go/test/data/data.txt`
	want := []string{
		"A B C D E F",
		"F  B  C  A B  E  F",
		"a bb ccc dd eeeee ffffff",
	}
	erg := ld.GetData(filename)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}
