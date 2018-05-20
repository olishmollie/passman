package passman

import (
	"io/ioutil"
	"path"
)

// Find finds, decrypts, and prints a password to the console
func Find(dir string) string {
	fname := path.Join(Root, dir)
	if !PswdExists(fname) {
		FatalError(nil, "cannot find pswd for "+dir)
	}
	ct, err := ioutil.ReadFile(fname)
	if err != nil {
		FatalError(err, "could not read pswd for "+dir)
	}
	k := getEncryptionKey()
	pswd, err := Decrypt(k, ct)
	if err != nil {
		FatalError(err, "bad encryption key for "+dir)
	}
	return string(pswd)
}
