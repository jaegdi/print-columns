// Package pc fomat and filter text columns
package main

import (
	"fmt"
	"os"

	ap "pc/internal/argparse"
	df "pc/internal/dataformat"
	ld "pc/internal/loaddata"
)

func init() {
	// ap.EvalFlags()
}

func main() {
	ap.CmdParams.Sep = " "
	ap.EvalFlags()
	// fmt.Println("Filename from CmdParams:", ap.CmdParams.Filename)
	// Load data from file and/or STDIN
	rawdata, err := ld.GetData(ap.CmdParams.Filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data: %v\n", err)
		os.Exit(1)
	}
	// fmt.Println("Raw data:", rawdata) // Added debug print statement for rawdata
	// Get the seperator for parsing the data input
	sep := []rune(ap.CmdParams.Sep)[0]
	//  parse the input data
	pdata := df.DataParse(rawdata, sep)
	// Format the parsed data and print out
	df.Format(pdata)
}
