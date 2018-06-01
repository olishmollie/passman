package passman

import (
	"path"
	"testing"
)

var root = "mock_store"
var keyfile = path.Join(root, ".key")
var infile = "mock_pswds.txt"

func TestFind(t *testing.T) {
	prefix1, pswd1 := "category/website/username1", "secret1"
	prefix2, pswd2 := "category/website/username2", "secret2"
	prefix3, pswd3 := "username3", "secret3"
	prefix4, pswd4 := "website/username4", "secret4"
	setupFind(root, keyfile)
	if res := Find(root, keyfile, prefix1, false); res != pswd1 {
		t.Errorf("Expected %s, got %s", pswd1, res)
	}
	if res := Find(root, keyfile, prefix2, false); res != pswd2 {
		t.Errorf("Expected %s, got %s", pswd2, res)
	}
	if res := Find(root, keyfile, prefix3, false); res != pswd3 {
		t.Errorf("Expected %s, got %s", pswd3, res)
	}
	if res := Find(root, keyfile, prefix4, false); res != pswd4 {
		t.Errorf("Expected %s, got %s", pswd4, res)
	}
	removeContentsOf(root, ".key")
}

func setupFind(root, keyfile string) {
	Import(root, keyfile, infile)
}
