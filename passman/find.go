package passman

import (
	"io/ioutil"
	"path"
)

// Find finds and decrypts a password in the prefix directory.
func Find(root, keyfile, prefix string) (p string, isDir bool) {
	fname := path.Join(root, prefix)
	if pswdExists(fname) {
		ct, err := ioutil.ReadFile(fname)
		if err != nil {
			FatalError(err, "could not read pswd for "+prefix)
		}
		pswd, err := decrypt(getEncryptionKey(keyfile), ct)
		if err != nil {
			FatalError(err, "bad encryption key for "+prefix)
		}
		p, isDir = string(pswd), false
		return
	} else if DirExists(fname) {
		p, isDir = fname, true
		return
	}
	FatalError(nil, "could not find pswd for "+prefix)
	return "", false
}
