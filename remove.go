package passman

import (
	"os"
	"path"
)

// Remove removes given password from storage
func Remove(p string) {
	dir := path.Join(Root, p)
	if DirExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			FatalError(err, "could not remove pswd for "+dir)
		}
	}
}
