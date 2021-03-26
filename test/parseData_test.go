package main

import (
	ap "pc/argparse"
	lp "pc/lineparse"
	ld "pc/loaddata"
	"reflect"
	"testing"
)

// Test GetLineString, line with fixed length columns seperated by two or more blanks
func TestParseDataFromFileLine1(t *testing.T) {
	filename := `/home/jaegdi/devel/go/pc-go/test/data.txt`
	want := lp.T_dataline{`A`, `B`, `C`, `D`, `E`, `F`}
	erg := lp.LineParse(ld.GetData(filename)[0], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}

func TestParseDataFromFileLine2(t *testing.T) {
	filename := `/home/jaegdi/devel/go/pc-go/test/data.txt`
	want := lp.T_dataline{`F`, `B`, `C`, `A_B`, `E`, `F`}
	erg := lp.LineParse(ld.GetData(filename)[1], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}
