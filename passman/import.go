package passman

import (
	"bufio"
	"os"
	"strings"
)

// Import reads pswds from a file and adds them to pswd store.
// File must be newline delimited directory and pswds (like what Dump outputs)
// e.g. Category/Website/username pswd
//      Category/NextWebsite/username pswd
func Import(root, keyfile, infile string) {
	f, err := os.Open(infile)
	if err != nil {
		FatalError(err, "could not open import file")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		d := strings.Split(t, " ")
		Add(root, keyfile, d[0], d[1])
	}
}
