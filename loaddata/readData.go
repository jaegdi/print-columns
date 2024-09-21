package pc

import (
	"bufio"
	"fmt"
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

// getFileData reads data from a file and returns it as a slice of strings.
// It takes the filename as an input parameter.
// If there's an error opening or reading the file, it logs a fatal error.
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

// GetData reads data from a file and/or stdin, depending on the provided parameters.
// It returns a slice of strings containing the read data.
func GetData(filename string) []string {
	data := []string{}

	// If a filename is provided in the command parameters, read from the file
	if ap.CmdParams.Filename != "" {
		fmt.Println("GetData filename:", filename) // Debug print: show the filename being processed
		data = getFileData(filename)               // Read data from the file
		fmt.Println("GetData data:", data)         // Debug print: show the data read from the file
	}

	// Check if there's input from stdin, and if so, append it to the data
	if checkStdin() {
		data = append(data, getStdinData()...) // Append data from stdin to the existing data
	}

	return data // Return the combined data from file and/or stdin
}
