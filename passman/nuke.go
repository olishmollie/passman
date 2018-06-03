package passman

import (
	"bytes"
	"os"
)

// Nuke removes passman folder and all its contents.
func Nuke(root string, force bool) {
	if !force {
		confirm := getln("CAUTION: This will remove all passwords. Continue? (y/N) ")
		confirm = bytes.TrimRight(confirm, "\n")
		if string(confirm) != "Y" && string(confirm) != "y" {
			Log("nuke aborted.")
			os.Exit(0)
		}
	}
	os.RemoveAll(root)
}
