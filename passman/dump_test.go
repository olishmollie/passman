package passman

import "testing"

var dumpFile = "pswds~"

func DumpTest(t *testing.T) {
	setupDump()
	Dump(root, keyfile, dumpFile)
	if !DirExists(dumpFile) {
		t.Errorf("Expected dump to create %s", dumpFile)
	}
	removeContentsOf(root, ".key")
}

func ImportTest(t *testing.T) {
	setupDump()
	Dump(root, keyfile, dumpFile)
	Import(root, keyfile, dumpFile)
	TestFind(t)
}

func setupDump() {
	prefix1, pswd1 := "category/website/username1", "secret1"
	prefix2, pswd2 := "category/website/username2", "secret2"
	prefix3, pswd3 := "username3", "secret3"
	prefix4, pswd4 := "website/username4", "secret4"
	Add(root, keyfile, prefix1, pswd1)
	Add(root, keyfile, prefix2, pswd2)
	Add(root, keyfile, prefix3, pswd3)
	Add(root, keyfile, prefix4, pswd4)
}
