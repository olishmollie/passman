package main

import (
	"os"
	"path"
	"testing"

	"github.com/olishmollie/passman/passman"
)

func setupDeleteTest() {
	passman.Add(root, keyfile, "delete/test", "pswd")
}

func TestDelete(t *testing.T) {
	setupDeleteTest()
	remove := new(test)
	remove.init(os.Stdin, os.Stdout, "delete", "delete/test")
	err := remove.run()
	if err != nil {
		t.Errorf("could not remove mock password\n%s", err)
	}

	if passman.FileExists(path.Join(root, "delete/test")) {
		t.Error("Remove did not delete prefix directory")
	}

}
