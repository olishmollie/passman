package passman

import (
	"path"
	"testing"
)

func TestRemove(t *testing.T) {
	setupRemove()
	prefixes := []string{"category/website/username1", "category/website/username2", "username3", "website/username4"}
	for i := 0; i < 2; i++ {
		Remove(root, prefixes[i])
		if DirExists(path.Join(root, prefixes[i])) {
			t.Errorf("Expected %s to have been removed", prefixes[i])
		}
		for j := i + 1; j < 3; j++ {
			if !DirExists(path.Join(root, prefixes[j])) {
				t.Errorf("Expeted %s to not have been removed", prefixes[j])
			}
		}
	}
	removeContentsOf(root, ".key")
}

func setupRemove() {
	Import(root, keyfile, infile)
}
