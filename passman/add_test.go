package passman

import (
	"io/ioutil"
	"path"
	"testing"
)

func TestAdd(t *testing.T) {
	var prefix, pswd = "category/website/username", "secret"
	addsDirectory(root, keyfile, prefix, pswd, t)
	storesCipher(root, keyfile, prefix, pswd, t)
	removeContentsOf(root, ".key")
}

func addsDirectory(root, keyfile, prefix, pswd string, t *testing.T) {
	Add(root, keyfile, prefix, pswd)
	if !DirExists(path.Join(root, prefix)) {
		t.Error("Expected Add() to create password directory")
	}
}

func storesCipher(root, keyfile, prefix, pswd string, t *testing.T) {
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
