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
	pdata := df.T_parsedData{}
	sep := []rune(ap.CmdParams.Colsep)[0]
	for _, v := range rawdata {
		l := df.LineParse(v, sep)
		pdata.Append(l)
	}
	df.Format(pdata)
}
