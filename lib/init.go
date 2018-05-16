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

	pub := path.Join(root, ".fpubkey")
	f, err := os.Create(pub)
	if err != nil {
		FatalError(err, "could not create key file")
	}

	k := GenKey(pswd)
	_, err = f.Write(k[:])
	if err != nil {
		FatalError(err, "could not write key to key file")
	}

	err = f.Close()
	if err != nil {
		FatalError(err, "could not close key file")
	}
}

func getUserKey() []byte {
	var pswd []byte
	for {
		fmt.Print("Please enter a passphrase. This will be used to generate a cipher key: ")
		pswd = Getln()
		fmt.Print("Confirm passphrase: ")
		pswd2 := Getln()
		if string(pswd) != string(pswd2) {
			fmt.Println("Passphrases don't match.")
		} else {
			break
		}
	}
	return pswd
}
