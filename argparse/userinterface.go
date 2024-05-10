package pc

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var CmdParams T_flags

// getArgsColNumbers collect all unknown parameters, if there are int values or ranges of int:int, as column numbers.
// ranges are supported - m:n, upwards 3:6 and downwards 6:3
func getArgsColNumbers() T_ColNumbers {
	var cn T_ColNumbers
	var error bool
	for _, val := range flag.Args() {
		if strings.Contains(val, ":") { // if a range
			res := strings.Split(val, ":") // split into fields
			if len(res) != 2 {
				error = true
				fmt.Println("ERROR: Wrong format for range: ", val)
			} else {
				n, errn := strconv.Atoi(res[0]) // convert  to int
				if errn != nil {
					error = true
					fmt.Println("ERROR: first part of range", val, "is not an integer. Range", val, "is ignored.")
				} else {
					m, errm := strconv.Atoi(res[1]) // convert  to int
					if errm != nil {
						error = true
						fmt.Println("ERROR: first part of range", val, "is not an integer. Range", val, "is ignored.")
					} else {
						if n < m {
							for i := n; i <= m; i++ { // upwards loop
								cn = append(cn, T_ColNum(i))
							}
						} else {
							for i := n; i >= m; i-- { // downwards loop
								cn = append(cn, T_ColNum(i))
							}
						}
					}
				}
			}
		} else {
			i, err := strconv.Atoi(val)
			if err == nil {
				cn = append(cn, T_ColNum(i))
			}
		}
	}
	if error {
		fmt.Println("program 'pc' is exited because of error in parameter!")
		os.Exit(1)
	}
	return cn
}

// fix_params disable CmdParams, that make no sense,when output to CSV or JSON.
func fix_params() {
	if CmdParams.Csv || CmdParams.Json {
		// CmdParams.Ts = false
		CmdParams.Fs = false
		CmdParams.Pp = false
	}
	CmdParams.Grouping = CmdParams.Gcol > 0
}

// EvalFlags evaluate all command line flags and set a struct with their values.
func EvalFlags() {
	flag.Usage = cmdManpage // helptext for parameters must be defined at function 'cmdParams'
	// Global Flags with values
	filenmPtr := flag.String("file", "", "Filename, read the text from this file")
	headerPtr := flag.String("header", "", "Headerline, if the text has no headers, you can define headers. They must be defined in the original order of the incoming text. Headers are left adjeusted, if they not start with a dash (-), then they right adjusted.")
	sepPtr := flag.String("sep", " ", "InputColumnSeperator, define the character to separate the columns, when parsing in, default=' '")
	colsepPtr := flag.String("colsep", "|", "ColumnSeperator, define the character to separate the columns, default='|'")
	filterPtr := flag.String("filter", "", "Filterpattern, process only lines where 'filter-string' is found")
	gcolnrPtr := flag.Int("gcol", 0, "GroupColumn, write a separator when the value in this column is different to the value in the previous line to group the values in this column. Number refers to the number of the output column")
	gcolvalPtr := flag.Bool("gcolval", false, "GroupColumnValues, Do not replace values in Groupcol by '' ")
	sortColPtr := flag.Int("sortcol", 0, "SortColumn, number of column, to sort for. Only one column ca be defined for sort.")
	colswPtr := flag.Int("w", 1, "colSepWidth, no of chars used to seperate output columns, default=1")
	// Boolean flags
	nfPtr := flag.Bool("nf", false, "no format, don't format the colums for common column width")
	nnPtr := flag.Bool("nn", false, "no format, don't format the numerical colums right adjusted")
	nhlPtr := flag.Bool("nhl", false, "no headline, The data contains no headline")
	tsPtr := flag.Bool("ts", false, "TitleSeparator, draws a separator line between first and second line of output")
	fsPtr := flag.Bool("fs", false, "FooterSeparator, draws a separator line between last and second last line of output")
	csPtr := flag.Bool("cs", false, "ColumnSeparator, draws a separator (default=|) between columns of output")
	ppPtr := flag.Bool("pp", false, "PrettyPrint, draw cell borders and all separators")
	rhPtr := flag.Bool("rh", false, "RemoveHeader, removes the first line")
	mbPtr := flag.Bool("mb", false, "MoreBlanks, more than one blank to split columns")
	numPtr := flag.Bool("num", false, "Num-bering, insert col numbers in the first line")
	csvPtr := flag.Bool("csv", false, "CSV, write output in CSV format")
	jsnPtr := flag.Bool("json", false, "JSON, write output in JSON format")
	jtcPtr := flag.Bool("jtc", false, "JSON, use first column as key")
	hlpPtr := flag.Bool("help", false, "Help, print help and exit")
	manPtr := flag.Bool("man", false, "Manual, print help and manual, then exit")
	verifyPtr := flag.Bool("v", false, "Verify, print parameter verirfy info")

	flag.Parse()

	// define map with all flags
	flags := T_flags{
		Filename:   string(*filenmPtr),
		Header:     string(*headerPtr),
		Sep:        string(*sepPtr),
		Colsep:     string(*colsepPtr),
		Filter:     string(*filterPtr),
		Gcol:       T_ColNum(*gcolnrPtr),
		GcolVal:    bool(*gcolvalPtr),
		SortCol:    T_ColNum(*sortColPtr),
		ColSepW:    int(*colswPtr),
		Nf:         bool(*nfPtr),
		Nn:         bool(*nnPtr),
		Nhl:        bool(*nhlPtr),
		Ts:         bool(*tsPtr),
		Fs:         bool(*fsPtr),
		Cs:         bool(*csPtr),
		Pp:         bool(*ppPtr),
		Rh:         bool(*rhPtr),
		Num:        bool(*numPtr),
		Csv:        bool(*csvPtr),
		Json:       bool(*jsnPtr),
		Jtc:        bool(*jtcPtr),
		Help:       bool(*hlpPtr),
		Manual:     bool(*manPtr),
		MoreBlanks: bool(*mbPtr),
		verify:     bool(*verifyPtr),
		Columns:    getArgsColNumbers(),
	}

	CmdParams = flags

	fix_params()

	if flags.Manual {
		cmdManpage()
		cmdExamples()
	}

	if flags.verify {
		println("\nCurrent values of parameters: ---------------------------------------")
		flags.Print()
	}
}
