package passman

import (
	"path"
	"testing"
)

func TestRemove(t *testing.T) {
	setup()

	Remove(test.root, mocks[0].prefix)

	if DirExists(path.Join(test.root, mocks[0].prefix)) {
		t.Error("Remove did not delete prefix directory")
	}

	teardown()
}
