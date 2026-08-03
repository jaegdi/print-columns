package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	"reflect"
	"testing"
)

// Test LineParse with data from file
func TestParseDataFromFileLine1(t *testing.T) {
	ap.CmdParams.Sep = " "
	ap.CmdParams.MoreBlanks = false

	data, err := readTestDataFile("data.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Line 1: "A B C D E F" - single space separated
	want := df.T_dataline{`A`, `B`, `C`, `D`, `E`, `F`}
	erg := df.LineParse(data[0], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse() = %q, want %q`, erg, want)
	}
}

func TestParseDataFromFileLine2(t *testing.T) {
	ap.CmdParams.Sep = " "
	ap.CmdParams.MoreBlanks = false

	data, err := readTestDataFile("data.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Line 2: "F  B  C  A B  E  F" - with MoreBlanks=false, all spaces are separators
	// Empty strings between double spaces are filtered out
	want := df.T_dataline{`F`, `B`, `C`, `A`, `B`, `E`, `F`}
	erg := df.LineParse(data[1], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse() = %q, want %q`, erg, want)
	}
}

func TestParseDataFromFileLine2WithMoreBlanks(t *testing.T) {
	ap.CmdParams.Sep = " "
	ap.CmdParams.MoreBlanks = true

	data, err := readTestDataFile("data.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Line 2: "F  B  C  A B  E  F" - with MoreBlanks=true, double spaces are separators
	// "A B" stays together as one field
	want := df.T_dataline{`F`, `B`, `C`, `A B`, `E`, `F`}
	erg := df.LineParse(data[1], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse() = %q, want %q`, erg, want)
	}
}

func TestAppendData1(t *testing.T) {
	want := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}}
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	ins := df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Append(ins)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf(`Append data %q with %q, got %q want match for %#q, nil`, tmp, ins, data, want)
	}
}

func TestInsertDataPos0(t *testing.T) {
	want := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	ins := df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Insert(ins, 0)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf(`Append data %q with %q, got %q want match for %#q, nil`, tmp, ins, data, want)
	}
}

func TestInsertDataPos1(t *testing.T) {
	want := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	ins := df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Insert(ins, 1)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf(`Append data %q with %q, got %q want match for %#q, nil`, tmp, ins, data, want)
	}
}

func TestInsertDataPos4(t *testing.T) {
	want := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	ins := df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Insert(ins, 4)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf(`Append data %q with %q, got %q want match for %#q, nil`, tmp, ins, data, want)
	}
}

func TestInsertDataPos5(t *testing.T) {
	want := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}}
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbbb`, `cccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccccc`, `ddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `dddddd`, `ee`, `f`},
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `eeeeee`, `f`}}
	ins := df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `ffffff`}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Insert(ins, 5)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf(`Append data %q with %q, got %q want match for %#q, nil`, tmp, ins, data, want)
	}
}
