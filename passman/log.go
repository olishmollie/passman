package passman

import (
	"fmt"
	"os"
)

// FatalError writes an error to stderr and exits.
func FatalError(err error, format string, args ...interface{}) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "fatal: "+format+"\n", args...)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "fatal: "+format+"\n"+err.Error()+"\n", args...)
	os.Exit(1)
}

// Log writes a message stdout.
func Log(msg string) {
	fmt.Printf("passman: %s\n", msg)
}
