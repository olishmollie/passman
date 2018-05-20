package passman

import (
	"fmt"
	"os"
)

// FatalError writes an error to stderr and exits
func FatalError(err error, msg string) {
	if err == nil {
		os.Stderr.WriteString("fatal: " + msg + "\n")
		os.Exit(1)
	}
	os.Stderr.WriteString("fatal: " + msg + "\n" + err.Error() + "\n")
	os.Exit(1)
}

// Log writes a message stdout
func Log(msg string) {
	fmt.Printf("passman: %s\n", msg)
}
