package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	ld "pc/loaddata"
	"reflect"
	"testing"
)

func TestMaxlenDataFromFileLine1(t *testing.T) {
	ap.CmdParams.MoreBlanks = true
	filename := `/home/jaegdi/devel/go/pc-go/test/data/data.txt`
	sep := []rune(ap.CmdParams.Sep)[0]
	data := ld.GetData(filename)
	want := df.T_maxlenghts{1, 2, 3, 3, 5, 6}
	d := df.DataParse(data, sep)
	erg := df.GetMaxLength(d)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}

func TestMaxlenData6to1(t *testing.T) {
	// ap.CmdParams.MoreBlanks = true
	data := df.T_parsedData{
		df.T_dataline{`aaaaaa`, `bbbbb`, `cccc`, `ddd`, `ee`, `f`}}
	want := df.T_maxlenghts{6, 5, 4, 3, 2, 1}
	erg := df.GetMaxLength(data)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
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
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}
