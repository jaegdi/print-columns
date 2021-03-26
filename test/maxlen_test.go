package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	lp "pc/lineparse"
	ld "pc/loaddata"
	"reflect"
	"testing"
)

func TestMaxlenDataFromFileLine1(t *testing.T) {
	filename := `/home/jaegdi/devel/go/pc-go/test/data.txt`
	want := lp.T_dataline{1, 1, 1, 1, 1, 1}
	data := lp.LineParse(ld.GetData(filename)[0], ' ')
	erg := df.GetMaxLength(data)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`GetData("%s") = %q, want match for %#q, nil`, ap.CmdParams.Filename, erg, want)
	}
}
