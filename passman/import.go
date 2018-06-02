package passman

import (
	"bufio"
	"os"
	"strings"
)

// Import unencrypts infile and imports passwords into store.
func Import(root, keyfile, infile string) {
	pswd := getUserPswd()
	key := hashPswd(pswd)
	decryptFile(key[:], infile)
	newKey := generateEncryptionKey()
	writeEncryptionKey(keyfile, newKey)
	addPswdFile(root, keyfile, infile)
}

func addPswdFile(root, keyfile, infile string) {
	f, err := os.Open(infile)
	if err != nil {
		FatalError(err, "could not open import file")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		d := strings.Split(t, " ")
		Add(root, keyfile, d[0], d[1])
	}
}
