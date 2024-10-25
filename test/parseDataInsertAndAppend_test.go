package main

import (
	"log"
	"reflect"
	"testing"

	ap "pc/internal/argparse"
	df "pc/internal/dataformat"
	ld "pc/internal/loaddata"
)

var filename string = "/home/dirk/devel/print-columns/test/data/data.txt"

// Test GetLineString, line with fixed length columns seperated by two or more blanks
func TestParseDataFromFileLine1(t *testing.T) {
	ap.CmdParams.MoreBlanks = true
	want := df.T_dataline{`A`, `B`, `9`, `C`, `D`, `E`, `F`, `3`}
	ap.CmdParams.Header = "A B 9 C D E F 3"
	data, err := ld.GetFileData(filename)
	if err != nil {
		t.Fatalf("GetData(\"%s\") returned error: %v", filename, err)
	}
	log.Println("want: ", want)
	log.Println("data: ", data)
	erg := df.LineParse(data[0], ' ')
	log.Println("erg: ", erg)
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, data[0], ' ', erg, want)
	}
}

func TestParseDataFromFileLine2(t *testing.T) {
	want := df.T_dataline{`F`, `B`, `5`, `C`, `A`, `B`, `E`, `F`}
	data, err := ld.GetFileData(filename)
	if err != nil {
		t.Fatalf("GetData(\"%s\") returned error: %v", filename, err)
	}
	erg := df.LineParse(data[1], ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, data[1], ' ', erg, want)
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
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "ee", "ffffff"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbbb", "cccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "dddddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "eeeeee", "f"}}
	data := df.T_parsedData{
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbbb", "cccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccccc", "ddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "dddddd", "ee", "f"},
		df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "eeeeee", "f"}}
	ins := df.T_dataline{"aaaaaa", "bbbbb", "cccc", "ddd", "ee", "ffffff"}
	tmp := make(df.T_parsedData, len(data))
	copy(tmp, data)
	data.Insert(ins, 0)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf("Append data \n%s\n with \n%s\n, got \n%s\n want match for \n%s\n", df.SprintData("", tmp), df.SprintData("", ins), df.SprintData("", data), df.SprintData("", want))
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
