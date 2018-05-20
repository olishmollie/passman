package passman

import (
	"os"
	"path"
)

// Add inserts a password into storage
func Add(p, data string) {

	dir, file := splitDir(p)

	newdir := path.Join(Root, dir)
	if !DirExists(newdir) {
		err := os.MkdirAll(newdir, 0755)
		if err != nil {
			FatalError(err, "could not create password store")
		}
	}

	f, err := os.Create(path.Join(Root, dir, file))
	if err != nil {
		FatalError(err, "could not create password")
	}

	ct, err := encrypt(getEncryptionKey(), []byte(data))
	if err != nil {
		FatalError(err, "could not encrypt password for "+dir)
	}

	_, err = f.Write(ct)
	if err != nil {
		FatalError(err, "could not write pswd for "+dir)
	}

	err = f.Close()
	if err != nil {
		FatalError(err, "could not close pswd for "+dir)
	}

}
