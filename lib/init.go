package lib

import (
	"fmt"
	"log"
	"os"
	"path"
)

// Init initializes passman by creating a storage directory and generating a cipher key
func Init() error {

	root := GetRootDir()

	fmt.Println("Welcome to Passman!")
	if !DirExists(root) {
		err := os.Mkdir(root, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	pswd := getUserKey()

	pub := path.Join(root, ".fpubkey")
	f, err := os.Create(pub)
	if err != nil {
		log.Fatal("error creating pub key file:", err)
	}

	k := GenKey(pswd)
	_, err = f.Write(k[:])
	if err != nil {
		log.Fatal("error writing to file:", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
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
