package pc

import (
	"bufio"
	"log"
	"os"
	ap "pc/argparse"
)

// get StdinData read data from STDIN
func getStdinData() []string {
	data := []string{}
	if checkStdin() {
		// fmt.Println("Something on STDIN")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}
	return data
}

// getFielData read data from file.
func getFileData(fname string) []string {
	data := []string{}
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Open File:"+fname, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Read File:"+fname, err)
	}
	return data
}

// GetDat read data from file, if param file is set and append data from STDIN, if there some data.
func GetData(filename string) []string {
	data := []string{}
	if ap.CmdParams.Filename != "" {
		data = getFileData(ap.CmdParams.Filename)
	}
	if checkStdin() {
		data = append(data, getStdinData()...)
	}
	return data
}
