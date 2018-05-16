package lib

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

// Init initializes passman by creating a storage directory and generating a cipher key
func Init() error {

	root := GetRootDir()
	if DirExists(root) {
		fmt.Println("Warning: password store is currently initialized.")
		fmt.Println("Reinitialization will result in the loss of your passwords.")
		fmt.Println("Do you wish to reinitialize? (y/N) ")
		reader := bufio.NewReader(os.Stdin)
		c, err := reader.ReadByte()
		if err != nil {
			FatalError(err, "could not read byte from stdin")
		}
		switch c {
		case 'y', 'Y':
			fmt.Println("Reinitializing...")
			err := os.RemoveAll(root)
			if err != nil {
				FatalError(err, "could not remove pswd store")
			}
		default:
			os.Exit(0)
		}
	}

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

	return nil
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
