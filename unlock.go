package passman

import "os"

// Unlock unencrypts the .passman.lock file and imports all passwords into the password store
func Unlock() {
	pswd := getUserPswd()
	key := hashPswd(pswd)
	DecryptFile(key[:], Lockfile)
	newKey := generateEncryptionKey()
	writeEncryptionKey(newKey)
	Import(Lockfile)
	os.Remove(Lockfile)
}
