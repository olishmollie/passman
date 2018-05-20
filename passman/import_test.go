package passman

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestImport(t *testing.T) {
	Import(root, keyfile, infile)
	wantedFiles := []string{"category", "username3", "website"}
	files := topLevelFiles(root)
	for i, el := range files {
		if el != wantedFiles[i] {
			t.Errorf("Expected %s to equal %s", el, wantedFiles[i])
		}
	}
	removeContentsOf(root, ".key")
}

func topLevelFiles(dir string) []string {
	allFiles, err := ioutil.ReadDir(root)
	if err != nil {
		FatalError(err, "could not read files from mock store")
	}
	files := []string{}
	for _, el := range allFiles {
		if strings.HasPrefix(el.Name(), ".") {
			continue
		}
		files = append(files, el.Name())
	}
	return files
}
