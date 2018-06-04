package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/olishmollie/passman/passman"
)

// var mocks = []mock{
// 	{"category/website/username1", "secret1"},
// 	{"category/website/username2", "secret2"},
// 	{"username3", "secret3"},
// 	{"website/username4", "secret4"},
// }

var findtests = []struct {
	in  string
	out string
}{
	{"category/website/username1", "secret1"},
	{"category/website/username2", "secret2"},
	{"username3", "secret3"},
	{"website/username4", "secret4"},
}

func TestFind(t *testing.T) {
	find := new(test)
	for _, tt := range findtests {
		passman.Add(root, keyfile, tt.in, tt.out)
		var b bytes.Buffer
		find.init(os.Stdin, &b, tt.in)
		err := find.run()
		if err != nil {
			t.Errorf("error running find command\n%s", err)
		}
		if b.String() != tt.out+"\n" {
			t.Error("did not match prefix to password")
		}
	}

}
