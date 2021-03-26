package pc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	ap "pc/argparse"
)

func getStdinData() []string {
	data := []string{}
	if checkStdin() {
		fmt.Println("Something on STDIN")
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

func GetData(filename string) []string {
	data := []string{}
	if filename == "" {
		filename = ap.CmdParams.Filename
	}
	if filename != "" {
		data = getFileData(filename)
	}
	if checkStdin() {
		data = append(data, getStdinData()...)
	}
	return data
}
