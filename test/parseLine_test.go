package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	"reflect"
	"testing"
)

// Reset params before each test
func resetParseParams() {
	ap.CmdParams.MoreBlanks = false
	ap.CmdParams.Sep = " "
}

// Test GetLineString, line with fixed length columns separated by two or more blanks
func TestGetLineStringFixedColumnsByBlanksOneCombinedTag(t *testing.T) {
	resetParseParams()
	ap.CmdParams.MoreBlanks = true

	line := "NAME       NAMESPACE                  DOCKER REF            ISTAG                 UPDATED"
	// With MoreBlanks=true, "DOCKER REF" should become one field because there's only one space between them
	// but multiple spaces separate the columns
	want := df.T_dataline{`NAME`, `NAMESPACE`, `DOCKER REF`, `ISTAG`, `UPDATED`}
	erg := df.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %q, want %q`, line, erg, want)
	}
}

func TestParseHeadlineFixedColumnsByBlanksMoreCombinedTags(t *testing.T) {
	resetParseParams()
	ap.CmdParams.MoreBlanks = true

	line := "NAME SPEC                 DOCKER REF                 UPDATED VAL"
	// With MoreBlanks, single spaces are kept together, multiple spaces are separators
	want := df.T_dataline{`NAME SPEC`, `DOCKER REF`, `UPDATED VAL`}
	erg := df.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %q, want %q`, line, erg, want)
	}
}

func TestParseDoubleQuoted(t *testing.T) {
	resetParseParams()

	line := `NAME "DOCKER REF" UPDATED`
	want := df.T_dataline{`NAME`, `"DOCKER REF"`, `UPDATED`}
	erg := df.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}

func TestParseSingleQuoted(t *testing.T) {
	resetParseParams()

	line := "NAME 'DOCKER REF' UPDATED"
	want := df.T_dataline{`NAME`, `'DOCKER REF'`, `UPDATED`}
	erg := df.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}

func TestParseVariableLength(t *testing.T) {
	resetParseParams()

	line := "NAME DOCKER REF UPDATED"
	want := df.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := df.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}

func TestParseComma(t *testing.T) {
	resetParseParams()

	line := "NAME,DOCKER,REF,UPDATED"
	want := df.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := df.LineParse(line, ',')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}

func TestParseCommaQuoted(t *testing.T) {
	resetParseParams()

	line := `NAME,"DOCKER,REF",UPDATED`
	want := df.T_dataline{`NAME`, `"DOCKER,REF"`, `UPDATED`}
	erg := df.LineParse(line, ',')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}

func TestParseTilde(t *testing.T) {
	resetParseParams()

	line := "NAME~DOCKER~REF~UPDATED"
	want := df.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := df.LineParse(line, '~')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s") = %v, want %v`, line, erg, want)
	}
}
