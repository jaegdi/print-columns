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
	ap.EvalFlags()
	rawdata := df.T_rawdata(ld.GetData(""))
	sep := []rune(ap.CmdParams.Sep)[0]
	pdata := df.DataParse(rawdata, sep)
	df.Format(pdata)
}
