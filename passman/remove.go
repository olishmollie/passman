package passman

import (
	"os"
	"path"
)

// Remove removes given password from storage
func Remove(root, prefix string) {
	dir := path.Join(root, prefix)
	if DirExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			FatalError(err, "could not remove pswd for "+dir)
		}
	}
}
