package passman

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
)

// Lock dumps the password store into an encrypted file, removes all passwords, and disables most operations.
// CAUTION: if you forget the password you use to lock passman, you will not be able to recover your passwords.
func Lock(root, keyfile, lockfile string) {
	Log("CAUTION - if you forget the password you use to lock passman, you will not be able to recover your passwords.")
	yes := getln("Do you wish to continue? (y/N) ")
	if string(yes) == "y" || string(yes) == "Y" {
		pswd := getUserPswd()
		k := hashPswd(pswd)
		f, err := os.OpenFile(lockfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			FatalError(err, "could not open lockfile")
		}
		Dump(root, root, keyfile, f)
		encryptFile(k[:], lockfile)
		Log("deleting password store...")
		removeContentsOf(root, ".passman.lock")
		Log("passman locked")
		err = f.Close()
		if err != nil {
			FatalError(err, "could not close lockfile")
		}
	} else {
		Log("lock operation aborted.")
	}
}

// Unlock unencrypts the .passman.lock file and imports all passwords into the password store.
func Unlock(root, keyfile, lockfile string) {
	pswd := getUserPswd()
	key := hashPswd(pswd)
	decryptFile(key[:], lockfile)
	newKey := generateEncryptionKey()
	writeEncryptionKey(keyfile, newKey)
	Import(root, keyfile, lockfile)
	os.Remove(lockfile)
}

func getUserPswd() []byte {
	var pswd []byte
	// TODO - limit number of times passphrase can be entered
	for {
		pswd = getln("Enter your passphrase: ")
		pswd2 := getln("Confirm passphrase: ")
		if string(pswd) != string(pswd2) {
			fmt.Println("Passphrases don't match.")
		} else {
			break
		}
	}
	return bytes.TrimRight(pswd, "\n")
}

func hashPswd(p []byte) [32]byte {
	return sha256.Sum256(p)
}
