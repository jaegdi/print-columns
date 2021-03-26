package pc

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var CmdParams T_flags

func cmdParams() {
	paramsTxt := `
    pc [options]  [column-numbers]

	options:        [-file=input-filename] [-header='col1-header col2-header coln-header'] [-colsep='|'] [-filter='string]
	                [-csv] [-json] [-ts] [-cs] [-rh] [-pp] [-num] [-h,-help] [-man]
	column-numbers: [[1] [2]...[n]]

	All the options must be set before the column-numbers.

    pc - (printColumns) formats text columns from stdin or file and print them as a ASCII table or CSV or JSON.
	     With the options the look of the output can be modified.
	     The columns for output can be selected by number.
         The input must have the same number of seperated columns in each line.
         The formated result is printed to stdout.

    Without parameters 'pc' print all columns from input in formated form

    The parameters:

        --file=filename                    read the text from this file,
                                           if there is also data from STDIN, this is added together


        --header='...'    Headerline,      if the text has no headers, you can define headers.
                                           They must be defined in the original order of the incoming text.
                                           Headers are left adjeusted, if they not start with a dash (-), then
                                           they right adjusted.
        --colsep='|'      ColumnSeparator  define the character to separate the columns, default='|'
        --filter='string' Filter lines,    process only lines where 'string' is found
        --gcol=colnum:    GroupCol         write a separator when the value in this column is different
                                           to the value in the previous line to group the values in this column
        -ts:              TitleSeparator,  draws a separator line between first and second line of output
        -fs:              FooterSeparator, draws a separator line between last and second last line of output
        -cs:              ColumnSeparator, draws a separator (default=|) between columns of output
        -pp:              PrettyPrint,     draw cell borders and all separators
        -rh:              RemoveHeader,    removes the first line
        -num:             Num-bering,      insert col numbers in the first line
        -csv:             CSV,             write output in CSV format
        -json:            JSON,            write output in JSON format

        1 2 .. n:         ColumnNumbers,   The number of the columns from the incoming text,
                                           that should printed out. To rearrange the columns
                                           the columns can given in the wanted order.

		-h -help:         Help,            print help and exit
        -man:             Manual,          print help and manual, then exit
	`
	fmt.Println(paramsTxt)
}

// cmdUsage print the man page
func cmdUsage() {
	usageText := `

    Examples:

    - A data-file 'data2.txt' has the following content:

        aaaaa bbbbbbbbbbbb cccccccc dd eeeeeee fffffffffff
        aaaaaaaaaaaaaaa bbbbbbbbbbbbb cc dddddddddddd eeeeaaaaaeeee ff
        aa bbbbb cc dd ee ffffffffff

    - To print the file with formated columns
        cmd: pc data2.txt
             pc --file=data2.txt
             or    cat data2.txt | pc
        result:
        aaaaa             bbbbbbbbbbbb    cccccccc   dd             eeeeeee         fffffffffff
        aaaaaaaaaaaaaaa   bbbbbbbbbbbbb   cc         dddddddddddd   eeeeaaaaaeeee   ff
        aa                bbbbb           cc         dd             ee              ffffffffff

    - To print that file formated with additional header:
        cmd: pc data2.txt header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 -col-4 col-5 -col-6'
        result:
        col-1             col-2           col-3             col-4   col-5                 col-6
        ---------------   -------------   --------   ------------   -------------   -----------
        aaaaa             bbbbbbbbbbbb    cccccccc   dd             eeeeeee         fffffffffff
        aaaaaaaaaaaaaaa   bbbbbbbbbbbbb   cc         dddddddddddd   eeeeaaaaaeeee   ff
        aa                bbbbb           cc         dd             ee              ffffffffff

    - To print that file formated with additional header and columnseparator:
        cmd: pc data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -cs
        result:
        col-1           | col-2         | col-3    | col-4        | col-5         | col-6
        --------------- | ------------- | -------- | ------------ | ------------- | -----------
        aaaaa           | bbbbbbbbbbbb  | cccccccc | dd           | eeeeeee       | fffffffffff
        aaaaaaaaaaaaaaa | bbbbbbbbbbbbb | cc       | dddddddddddd | eeeeaaaaaeeee | ff
        aa              | bbbbb         | cc       | dd           | ee            | ffffffffff

    - To print that file formated with additional header and prettyprint:
        cmd: pc data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -pp
        result:
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |
        | col-1           | col-2         | col-3    | col-4        | col-5         | col-6       |
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |
        | aaaaa           | bbbbbbbbbbbb  | cccccccc | dd           | eeeeeee       | fffffffffff |
        | aaaaaaaaaaaaaaa | bbbbbbbbbbbbb | cc       | dddddddddddd | eeeeaaaaaeeee | ff          |
        | aa              | bbbbb         | cc       | dd           | ee            | ffffffffff  |
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |

    - To print col2 and col5 with additional headers in reverse order
        cmd: pc data2.txt --header='col-2 col-5 ' 5 2
             or cat data2.txt | pc --header='col-2 col-5 ' 5 2
        result:
        col-5           col-2
        -------------   -------------
        eeeeeee         bbbbbbbbbbbb
        eeeeaaaaaeeee   bbbbbbbbbbbbb
        ee              bbbbb

    - Format the output of a command
        cmd: oc get pod -o wide --all-namespaces | pc -ts -cs --gcol=2   8 1 2 5 6
        result:
        NODE                      | NAMESPACE             | NAME                                               | RESTARTS | AGE
        ------------------------- | --------------------- | -------------------------------------------------- | -------- | ---
        int-apc0-wrk-v08.sf-rz.de | app-monitoring        | prometheus-prometheus-0                            | 11       | 1d
        int-apc0-wrk-v10.sf-rz.de | br-test               | rsync-container-1-trkwt                            | 1        | 27d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | ldapgroupsync-1583331300-86bg8                     | 0        | 22d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | ldapgroupsync-1583334900-fsh48                     | 0        | 22d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | prune-builds-1585239000-lrncj                      | 0        | 1h
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | prune-deployments-1585242300-vr22s                 | 0        | 24m
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | registry-image-pruning-1585235220-prbj7            | 0        | 2h
        int-apc0-inf-v03.sf-rz.de | default               | docker-registry-5-bxk5x                            | 0        | 27d
        int-apc0-mst-v00.sf-rz.de | default               | registry-console-7-sj72f                           | 0        | 8d

    - Filter the output of a command and convert to json
        cmd:  oc get pod -o wide --all-namespaces |pc -ts -cs -json --filter='wrk-v01'   8 1 2 5 6
        result:
        {
            "data": [
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "fpc-fa2",
                    "datenkopie-zulieferung-46-46dhb",
                    "1",
                    "8h"
                ],
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "fpc-int1",
                    "datenkopie-zulieferung-64-pdp5r",
                    "1",
                    "8h"
                ],
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "openshift-logging",
                    "logging-fluentd-6bg5h",
                    "3",
                    "23d"
                ]
            ]
        }
	`
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func getArgsColNumbers() T_ColNumbers {
	var cn T_ColNumbers
	for _, val := range flag.Args() {
		i, err := strconv.Atoi(val)
		if err == nil {
			cn = append(cn, T_ColNum(i))
		}
	}
	return cn
}

func fix_params() {
	if CmdParams.Csv || CmdParams.Json {
		CmdParams.Ts = false
		CmdParams.Fs = false
		CmdParams.Pp = false
	}
	CmdParams.Grouping = CmdParams.Gcol > 0
}

// EvalFlags evaluate all command line flags and set a struct with their values
func EvalFlags() {
	flag.Usage = cmdParams
	// Global Flags
	filenmPtr := flag.String("file", "", "Filename,        read the text from this file")
	headerPtr := flag.String("header", "", "Headerline,      if the text has no headers, you can define headers. They must be defined in the original order of the incoming text. Headers are left adjeusted, if they not start with a dash (-), then they right adjusted.")
	colsepPtr := flag.String("colsep", " ", "ColumnSeperator, define the character to separate the columns, default='|'")
	filterPtr := flag.String("filter", "", "Filterpattern,   process only lines where 'filter-string' is found")
	gcolnrPtr := flag.Int("gcol", 0, "GroupColumn,     write a separator when the value in this column is different to the value in the previous line to group the values in this column")

	tsPtr := flag.Bool("ts", false, "TitleSeparator,  draws a separator line between first and second line of output")
	fsPtr := flag.Bool("fs", false, "FooterSeparator, draws a separator line between last and second last line of output")
	csPtr := flag.Bool("cs", false, "ColumnSeparator, draws a separator (default=|) between columns of output")
	ppPtr := flag.Bool("pp", false, "PrettyPrint,     draw cell borders and all separators")
	rhPtr := flag.Bool("rh", false, "RemoveHeader,    removes the first line")
	numPtr := flag.Bool("num", false, "Num-bering,      insert col numbers in the first line")
	csvPtr := flag.Bool("csv", false, "CSV,             write output in CSV format")
	jsnPtr := flag.Bool("json", false, "JSON,            write output in JSON format")
	hlpPtr := flag.Bool("help", false, "Help,            print help and exit")
	manPtr := flag.Bool("man", false, "Manual,          print help and manual, then exit")

	flag.Parse()

	// define map with all flags
	flags := T_flags{
		Filename: string(*filenmPtr),
		Header:   string(*headerPtr),
		Colsep:   string(*colsepPtr),
		Filter:   string(*filterPtr),
		Gcol:     T_ColNum(*gcolnrPtr),
		Ts:       bool(*tsPtr),
		Fs:       bool(*fsPtr),
		Cs:       bool(*csPtr),
		Pp:       bool(*ppPtr),
		Rh:       bool(*rhPtr),
		Num:      bool(*numPtr),
		Csv:      bool(*csvPtr),
		Json:     bool(*jsnPtr),
		Help:     bool(*hlpPtr),
		Manual:   bool(*manPtr),
		Columns:  getArgsColNumbers(),
	}

	CmdParams = flags

	fix_params()

	if flags.Manual {
		cmdParams()
		cmdUsage()
		fmt.Println("\nCurrent values of parameters: ---------------------------------------")
		flags.Print()
	}
}
