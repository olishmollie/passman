package passman

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// Print displays all the password prefixes in s if s is a prefix, otherwise prints password.
func Print(s string, offset int) {
	const lbar = "|\u2014\u2014 "
	files, err := ioutil.ReadDir(s)
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
			p := path.Join(s, n)
			Print(p, offset+2)
		}
	}
}
