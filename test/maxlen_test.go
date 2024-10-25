package main

import (
	"reflect"
	"testing"

	ap "pc/internal/argparse"
	df "pc/internal/dataformat"
	ld "pc/internal/loaddata"
)

func TestMaxlenDataFromFileLine1(t *testing.T) {
	filename := "./data/data.txt"
	ap.CmdParams.Filename = filename
	ap.CmdParams.MoreBlanks = true
	// fmt.Println("\nCurrent working directory:", os.Getenv("PWD"))
	sep := ' '
	// fmt.Println("TestMaxlenDataFromFileLine1 filename:", filename)
	data, err := ld.GetData(filename)
	if err != nil {
		t.Fatalf("GetData(\"%s\") returned error: %v", filename, err)
	}
	// PrintData("Data read from file", data)
	want := df.T_maxlenghts{1, 2, 1, 3, 2, 5, 6, 1}
	// PrintData("want:", want)
	d := df.DataParse(data, sep)
	// PrintData("d:", d)
	erg := df.GetMaxLength(d)
	// PrintData("erg:", erg)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %s, want match for %s`, filename, df.SprintData("erg:", erg), df.SprintData("want:", want))
	}
}

func TestMaxlenData6to1(t *testing.T) {
	// ap.CmdParams.MoreBlanks = true
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`}}
	want := df.T_maxlenghts{6, 5, 4, 3, 2, 1}
	erg := df.GetMaxLength(data)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q`, ap.CmdParams.Filename, erg, want)
	}
}

func TestMaxlenData6to6(t *testing.T) {
	// ap.CmdParams.MoreBlanks = true
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
		t.Fatalf(`GetData("%s") = %q, want match for %#q`, ap.CmdParams.Filename, erg, want)
	}
}
