package pc

import (
	"os"
)

/**
 * @brief Checks if there is data available on standard input (stdin).
 *
 * Determines whether data is available on standard input by checking the file mode of stdin. Returns true if data is available, false otherwise.
 * @return bool True if data is available on stdin, false otherwise.
 */
// checkStdin returns true if there are data on STDIN
func checkStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	} else {
		return true
	}
}
