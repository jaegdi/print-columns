package pc

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type T_ColNum int
type T_ColNumbers []T_ColNum

type T_flags struct {
	Filename   string
	Header     string
	Sep        string
	Colsep     string
	Filter     string
	Gcol       T_ColNum
	GcolVal    bool
	SortCol    T_ColNum
	ColSepW    int
	Nf         bool
	Nn         bool
	Nhl        bool
	Ts         bool
	Fs         bool
	Cs         bool
	Pp         bool
	Rh         bool
	Num        bool
	Csv        bool
	Json       bool
	Jtc        bool
	Help       bool
	Manual     bool
	Grouping   bool
	MoreBlanks bool
	verify     bool
	Columns    T_ColNumbers
}

func (t T_flags) Print() {
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		fmt.Fprintf(os.Stderr, "  -%-15s: %v\n", strings.ToLower(v.Type().Field(i).Name), v.Field(i))
	}
}
