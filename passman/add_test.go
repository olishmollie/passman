package passman

import (
	"io/ioutil"
	"path"
	"testing"
)

var root = "mock_store"
var keyfile = path.Join(root, ".key")
var prefix, pswd = "category/website/username", "secret"

func TestAdd(t *testing.T) {
	addsDirectory(t)
	storesCipher(t)
	removeContentsOf(root, ".key")
}

func addsDirectory(t *testing.T) {
	Add(root, keyfile, prefix, pswd)
	if !DirExists(path.Join(root, prefix)) {
		t.Error("Expected Add() to create password directory")
	}
}

func storesCipher(t *testing.T) {
	cipher, err := ioutil.ReadFile(path.Join(root, prefix))
	if err != nil {
		t.Error("error reading mock password")
	}
	data, err := decrypt(getEncryptionKey(keyfile), cipher)
	if err != nil {
		t.Error("error decrypting mock data")
	}
	if string(data) != pswd {
		t.Errorf("Expected %s to equal %s\n", string(data), pswd)
	}
}
