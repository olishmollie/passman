package passman

import (
	"io/ioutil"
	"path"
	"testing"
)

func TestAdd(t *testing.T) {
	setup()

	mock := mock{"foo/bar/baz", "boom"}
	Add(test.root, test.keyfile, mock.prefix, mock.pswd)

	pswd, err := ioutil.ReadFile(path.Join(test.root, mock.prefix))
	if err != nil {
		FatalError(err, "could not read file in TestAdd")
	}

	if string(pswd) == mock.pswd {
		t.Error("password was not encrypted")
	}

	teardown()
}
