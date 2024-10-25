package pc

import (
	"bufio"
	"fmt"
	"io"
	"os"

	ap "pc/internal/argparse"
)

/**
 * @brief Reads data from standard input (stdin).
 *
 * Efficiently reads data from standard input using bufio.Scanner, pre-allocating a slice to improve performance.  Returns an error if reading from stdin fails.
 * @return ([]string, error) A slice of strings representing the data read from stdin and an error, if any.
 */
// getStdinData reads data from STDIN, pre-allocating a slice to improve efficiency.
func getStdinData() ([]string, error) {
	const initialCapacity = 10000 // Adjust as needed based on typical input size
	data := make([]string, 0, initialCapacity)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}
	return data, nil
}

/**
 * @brief Reads data from a file.
 *
 * Efficiently reads data from a file using bufio.NewReader, pre-allocating a slice to improve performance.  Includes robust error handling for file operations.
 * @param fname The path to the file.
 * @return ([]string, error) A slice of strings representing the data read from the file and an error, if any.
 */
// getFileData reads data from a file, pre-allocating a slice to improve efficiency.  Includes improved error handling.
func GetFileData(fname string) ([]string, error) {
	const initialCapacity = 10000 // Adjust as needed based on typical input size
	data := make([]string, 0, initialCapacity)
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", fname, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file '%s': %w", fname, err)
		}
		data = append(data, line[:len(line)-1]) // Remove trailing newline
	}
	return data, nil
}

/**
 * @brief Reads data from a file and/or standard input (stdin).
 *
 * Reads data from a file if a filename is specified in the command-line parameters, and appends data from standard input if available.  Returns the combined data. Includes robust error handling.
 * @param filename The path to the file (optional).
 * @return ([]string, error) A slice of strings containing the combined data from file and/or stdin, and an error if any occurred during reading.
 */
// GetData reads data from a file and/or stdin, pre-allocating slices for efficiency.  Includes improved error handling.
func GetData(filename string) ([]string, error) {
	data := make([]string, 0)

	// If a filename is provided in the command parameters, read from the file
	if ap.CmdParams.Filename != "" {
		// fmt.Println("GetData filename:", filename) // Debug print: show the filename being processed
		fileData, err := GetFileData(filename) // Read data from the file
		if err != nil {
			return nil, err
		}
		data = append(data, fileData...)
		// fmt.Println("GetData data:", data) // Debug print: show the data read from the file
	}

	// Check if there's input from stdin, and if so, append it to the data
	if checkStdin() {
		stdinData, err := getStdinData() // Append data from stdin to the existing data
		if err != nil {
			return nil, err
		}
		data = append(data, stdinData...)
	}

	return data, nil // Return the combined data from file and/or stdin
}
