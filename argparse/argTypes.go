package pc

import (
	"fmt"
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
	Ts         bool
	Fs         bool
	Cs         bool
	Pp         bool
	Rh         bool
	Num        bool
	Csv        bool
	Json       bool
	Help       bool
	Manual     bool
	Grouping   bool
	MoreBlanks bool
	Columns    T_ColNumbers
}

func (t T_flags) Print() {
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("  -%-15s: %v\n", strings.ToLower(v.Type().Field(i).Name), v.Field(i))
	}
}
