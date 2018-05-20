package passman

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func getln(msg string) []byte {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadBytes('\n')
	if err != nil {
		FatalError(err, "could not read from stdin")
	}
	out := bytes.TrimRight(in, "\n")
	return out
}
