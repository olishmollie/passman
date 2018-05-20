package passman

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Dump writes unecrypted passwords to outfile.
func Dump(start, root, keyfile string, outfile *os.File) {
	files, err := ioutil.ReadDir(start)
	if err != nil {
		FatalError(err, "could not read files from password store")
	}
	for _, f := range files {
		n := f.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		if f.IsDir() {
			p := path.Join(start, n)
			Dump(p, root, keyfile, outfile)
		} else {
			p := path.Join(start, n)
			c, err := ioutil.ReadFile(p)
			if err != nil {
				FatalError(err, "could not read pswd for "+p)
			}
			t, err := decrypt(getEncryptionKey(keyfile), c)
			w := strings.TrimPrefix(p, root+"/") + " " + string(t) + "\n"
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
