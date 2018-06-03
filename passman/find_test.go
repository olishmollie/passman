package passman

import (
	"testing"
)

func TestFind(t *testing.T) {
	setup()
	for _, mock := range mocks {
		if pswd, _ := Find(test.root, test.keyfile, mock.prefix); pswd != mock.pswd {
			t.Errorf("expected %s, got %s", mock.pswd, pswd)
		}
	}
	teardown()
}
