// Package pc fomat and filter text columns
package main

import (
	ap "pc/argparse"
	df "pc/dataformat"
	ld "pc/loaddata"
)

func init() {
	// ap.EvalFlags()
}

func main() {
	ap.CmdParams.Sep = " "
	ap.EvalFlags()
	// Load data from file and/or STDIN
	rawdata := df.T_rawdata(ld.GetData(""))
	// Get the seperator for parsing the data input
	sep := []rune(ap.CmdParams.Sep)[0]
	//  parse the input data
	pdata := df.DataParse(rawdata, sep)
	// Format the parsed data and print out
	df.Format(pdata)
}
