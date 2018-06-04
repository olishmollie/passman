package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var addtests = [][]string{
	{"add", "foo/bar/baz", "boom"},
	{"add", "foo/blah", "foobar"},
	{"add", "foo/bar/bing", "bang"},
}

func TestAdd(t *testing.T) {

	for _, tt := range addtests {
		add := new(test)
		add.init(os.Stdin, os.Stdout, tt...)
		err := add.run()
		if err != nil {
			t.Errorf("could not run add command\n%s", err)
		}
		pswd, err := ioutil.ReadFile(path.Join(root, tt[1]))
		if err != nil {
			t.Errorf("could not find mock password\n%s", err)
		}

		if string(pswd) == tt[2] {
			t.Error("password was not encrypted")
		}
	}
}
