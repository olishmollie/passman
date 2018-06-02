package passman

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// Print displays all the password prefixes in dir.
func Print(dir string, offset int) {
	const lbar = "|\u2014\u2014 "
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not read files in password store")
	}
	for _, f := range files {
		var spaces []rune
		for i := 0; i < offset; i++ {
			spaces = append(spaces, ' ')
			spaces = append(spaces, ' ')
		}
		n := f.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		fmt.Println(string(spaces) + lbar + n)
		if f.IsDir() {
			p := path.Join(dir, n)
			Print(p, offset+2)
		}
	}
}
