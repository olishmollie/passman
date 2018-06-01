package passman

import (
	"io/ioutil"
	"path"

	"github.com/atotto/clipboard"
)

// Find finds, decrypts, and prints a password to the console
func Find(root, keyfile, prefix string, copy bool) string {
	fname := path.Join(root, prefix)
	if !pswdExists(fname) {
		FatalError(nil, "cannot find pswd for "+prefix)
	}
	ct, err := ioutil.ReadFile(fname)
	if err != nil {
		FatalError(err, "could not read pswd for "+prefix)
	}
	pswd, err := decrypt(getEncryptionKey(keyfile), ct)
	if err != nil {
		FatalError(err, "bad encryption key for "+prefix)
	}
	if copy {
		clipboard.WriteAll(string(pswd))
	}
	return string(pswd)
}
