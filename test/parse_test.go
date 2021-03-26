package main

import (
	lp "pc/lineparse"
	"reflect"
	"regexp"
	"testing"
)

// Test GetLineString, line with fixed length columns seperated by two or more blanks
func TestGetLineStringFixedColumnsByBlanksOneCombinedTag(t *testing.T) {
	line := "NAME       NAMESPACE                  DOCKER REF            ISTAG                 UPDATED"
	want := regexp.MustCompile(`NAME NAMESPACE DOCKER_REF ISTAG UPDATED`)
	msg := lp.GetLineString(line)
	if !want.MatchString(msg) {
		t.Fatalf(`Hello("%s") = %q, want match for %#q, nil`, line, msg, want)
	}
}

func TestParseHeadlineFixedColumnsByBlanksMoreCombinedTags(t *testing.T) {
	line := "NAME SPEC                 DOCKER REF                 UPDATED VAL"
	want := lp.T_dataline{`NAME_SPEC`, `DOCKER_REF`, `UPDATED_VAL`}
	erg := lp.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`Hello("%s") = %q, want match for %#q, nil`, line, erg, want)
	}
}

func TestParseDoubleQuoted(t *testing.T) {
	line := "NAME \"DOCKER REF\" UPDATED"
	want := lp.T_dataline{`NAME`, `"DOCKER REF"`, `UPDATED`}
	erg := lp.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}

func TestParseSingleQuoted(t *testing.T) {
	line := "NAME 'DOCKER REF' UPDATED"
	want := lp.T_dataline{`NAME`, `'DOCKER REF'`, `UPDATED`}
	erg := lp.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}

func TestParseVariableLength(t *testing.T) {
	line := "NAME DOCKER REF UPDATED"
	want := lp.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := lp.LineParse(line, ' ')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}

func TestParseComma(t *testing.T) {
	line := "NAME,DOCKER,REF,UPDATED"
	want := lp.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := lp.LineParse(line, ',')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}

func TestParseCommaQuoted(t *testing.T) {
	line := `NAME,"DOCKER,REF",UPDATED`
	want := lp.T_dataline{`NAME`, `"DOCKER,REF"`, `UPDATED`}
	erg := lp.LineParse(line, ',')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}

func TestParseTilde(t *testing.T) {
	line := "NAME~DOCKER~REF~UPDATED"
	want := lp.T_dataline{`NAME`, `DOCKER`, `REF`, `UPDATED`}
	erg := lp.LineParse(line, '~')
	if !reflect.DeepEqual(erg, want) {
		t.Fatalf(`LineParse("%s",'%v') = %v, want match for %v`, line, ' ', erg, want)
	}
}
