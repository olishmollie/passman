package passman

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func getln(in io.Reader, msg string) []byte {
	fmt.Print(msg)
	reader := bufio.NewReader(in)
	input, err := reader.ReadBytes('\n')
	if err != nil {
		FatalError(err, "could not read from stdin")
	}
	out := bytes.TrimRight(input, "\n")
	return out
}

func copyFile(f string) string {

	tmp, err := os.OpenFile(f+".tmp", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		FatalError(err, "could not open tmp file for copying")
	}
	defer tmp.Close()

	data, err := ioutil.ReadFile(f)
	if err != nil {
		FatalError(err, "could not read file for copying")
	}

	_, err = tmp.Write(data)
	if err != nil {
		FatalError(err, "could not write to file for copying")
	}

	return tmp.Name()

}
