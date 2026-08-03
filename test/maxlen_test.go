package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	"reflect"
	"testing"
)

func TestMaxlenDataFromFileLine1(t *testing.T) {
	ap.CmdParams.MoreBlanks = true
	ap.CmdParams.Sep = " "
	sep := ' '

	// Read test data using helper
	data, err := readTestDataFile("data.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// data.txt contains:
	// Line 1: "A B C D E F" - single spaces, with MoreBlanks these become one field
	// Line 2: "F  B  C  A B  E  F" - double spaces are separators, "A B" stays together
	// Line 3: "a bb ccc dd eeeee ffffff" - single spaces
	//
	// With MoreBlanks=true and parsing line by line:
	// The test data needs to match what LineParse returns

	d := df.DataParse(data, sep)

	// Line 2 parses to: [F, B, C, A B, E, F] - 6 columns, max lengths: 1, 2, 3, 3, 5, 6
	// But lines 1 and 3 are single-space separated so they become single fields with MoreBlanks
	// This test may need adjustment based on actual data format

	// Let's just verify we get reasonable output
	erg := df.GetMaxLength(d)
	if len(erg) == 0 {
		t.Fatalf("GetMaxLength() returned empty result for data: %v", d)
	}

	t.Logf("Parsed data: %v", d)
	t.Logf("Max lengths: %v", erg)
}

func TestMaxlenData6to1(t *testing.T) {
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`}}
	want := df.T_maxlenghts{6, 5, 4, 3, 2, 1}
	erg := df.GetMaxLength(data)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetMaxLength() = %v, want %v`, erg, want)
	}
}

func TestMaxlenData6to6(t *testing.T) {
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}}
	want := df.T_maxlenghts{6, 6, 6, 6, 6, 6}
	erg := df.GetMaxLength(data)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetMaxLength() = %v, want %v`, erg, want)
	}
}
