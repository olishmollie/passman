package lib

import (
	"fmt"
	"os"
	"path"
)

// Init initializes passman by creating a storage directory and generating a cipher key
func Init() {

	root := GetRootDir()

	fmt.Println("Welcome to Passman!")
	if !DirExists(root) {
		err := os.Mkdir(root, 0755)
		if err != nil {
			FatalError(err, "could not create pswd store")
		}
	}

	pswd := getUserKey()
	writeUserKey(pswd)
}

func getUserKey() []byte {
	var pswd []byte
	// TODO - limit number of times passphrase can be entered
	for {
		pswd = Getln("Enter your passphrase: ")
		pswd2 := Getln("Confirm passphrase: ")
		if string(pswd) != string(pswd2) {
			fmt.Println("Passphrases don't match.")
		} else {
			break
		}
	}
	return pswd
}

func writeUserKey(pswd []byte) {
	root := GetRootDir()
	pub := path.Join(root, ".fpubkey")
	f, err := os.Create(pub)
	if err != nil {
		FatalError(err, "could not create key file")
	}

	k := genKey(pswd)
	_, err = f.Write(k[:])
	if err != nil {
		FatalError(err, "could not write key to key file")
	}

	err = f.Close()
	if err != nil {
		FatalError(err, "could not close key file")
	}
}
