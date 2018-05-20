package passman

import (
	"os"
	"path"
)

// Lock dumps the password store into an encrypted file and removes all passwords and .key
// CAUTION: if you forget the password you used to generate the encryption key, you will not
// be able to unencrypt your passwords.
func Lock() {
	Log("CAUTION - if you forget the password you use to lock passman, you will not be able to unencrypt your passwords.")
	yes := Getln("Do you wish to continue? (y/N) ")
	if string(yes) == "y" || string(yes) == "Y" {
		pswd := getUserPswd()
		k := hashPswd(pswd)
		lockfile := path.Join(Root, ".passman.lock")
		f, err := os.OpenFile(lockfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			FatalError(err, "could not open lockfile")
		}
		Dump(Root, f)
		EncryptFile(k[:], lockfile)
		Log("deleting password store...")
		RemoveContents(Root, ".passman.lock")
		Log("passman locked")
		err = f.Close()
		if err != nil {
			FatalError(err, "could not close lockfile")
		}
	} else {
		Log("lock operation aborted.")
	}
}
