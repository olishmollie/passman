package passman

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Dump writes unecrypted passwords to outfile.
func Dump(dir string, outfile *os.File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not read files from password store")
	}
	for _, f := range files {
		n := f.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		if f.IsDir() {
			p := path.Join(dir, n)
			Dump(p, outfile)
		} else {
			p := path.Join(dir, n)
			c, err := ioutil.ReadFile(p)
			if err != nil {
				FatalError(err, "could not read pswd for "+p)
			}
			t, err := decrypt(getEncryptionKey(), c)
			w := strings.TrimPrefix(p, Root+"/") + " " + string(t) + "\n"
			if err != nil {
				FatalError(err, "could not decrypt pswd for "+w)
			}
			_, err = outfile.WriteString(w)
			if err != nil {
				FatalError(err, "could not write passwords to dump file")
			}
		}
	}
}
