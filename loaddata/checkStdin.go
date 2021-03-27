package pc

import (
	"os"
)

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
