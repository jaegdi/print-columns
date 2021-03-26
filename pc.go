package main

import (
	ap "pc/argparse"
	lp "pc/lineparse"
	ld "pc/loaddata"
)

func init() {
	// ap.EvalFlags()
}

func main() {
	ap.EvalFlags()
	rawdata := lp.T_rawdata(ld.GetData(""))
	pdata := lp.T_parsedData{}
	sep := []rune(ap.CmdParams.Colsep)[0]
	for _, v := range rawdata {
		l := lp.LineParse(v, sep)
		pdata = pdata.Append(l)
	}
	pdata.Print()
}
