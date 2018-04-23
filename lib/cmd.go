package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Print prints a tree of the password store
func Print(dir string, offset int) {
	const lbar = "|\u2014\u2014 "
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("error reading files in password store | ", err)
	}
	for _, f := range files {
		var spaces []rune
		for i := 0; i < offset; i++ {
			spaces = append(spaces, ' ')
			spaces = append(spaces, ' ')
		}
		n := f.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		fmt.Println(string(spaces) + lbar + n)
		if f.IsDir() {
			p := path.Join(dir, n)
			Print(p, offset+2)
		}
	}
}

// Find finds, decrypts, and prints a password to the console
func Find(dir string) string {
	root := GetRootDir()
	fname := path.Join(root, dir)
	if !PswdExists(fname) {
		fmt.Println("error: that password doesn't exist.")
		os.Exit(1)
	}
	ct, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	k := GetKey()
	pswd, err := Decrypt(k, ct)
	if err != nil {
		log.Fatal(err)
	}
	return string(pswd)
}

// Add inserts a password into storage
func Add(p, data string) {

	root := GetRootDir()
	dir, file := SplitDir(p)

	newdir := path.Join(root, dir)
	if !DirExists(newdir) {
		err := os.MkdirAll(newdir, 0755)
		if err != nil {
			log.Fatal("error creating password directory | ", err)
		}
	}

	if PswdExists(path.Join(root, p)) {
		fmt.Println("error: that password already exists")
		os.Exit(1)
	}

	f, err := os.Create(path.Join(root, dir, file))
	if err != nil {
		log.Fatal("error making password file | ", err)
	}

	k := GetKey()
	ct, err := Encrypt(k, []byte(data))
	if err != nil {
		log.Fatal("error encrypting password data | ", err)
	}

	_, err = f.Write(ct)
	if err != nil {
		log.Fatal("error writing to file | ", err)
	}

}

// Remove removes given password from storage
func Remove(p string) {

	root := GetRootDir()

	dir := path.Join(root, p)
	if DirExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Fatal("error removing password | ", err)
		}
	}
}

// Edit decrypts and opens password file in editor defined by $VISUAL
func Edit(p string) {

	root := GetRootDir()
	f := path.Join(root, p)
	if !PswdExists(f) {
		fmt.Println("error: that password doesn't exist")
		os.Exit(1)
	}

	k := GetKey()

	// Read encrypted password from file
	ciphertext, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("unable to read password file | ", err)
	}

	// Decrypt password
	plaintext, err := Decrypt(k, ciphertext)
	if err != nil {
		log.Fatal("unable to decode password file | ", err)
	}

	// Write plaintext to password file
	ioutil.WriteFile(f, plaintext, 0775)

	// Open in editor
	cmd := exec.Command(os.ExpandEnv("$VISUAL"), f)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("unable to run edit command | ", err)
	}

	// Read edited text from file
	plaintext, err = ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("unable to read from password file after editing | ", err)
	}

	// For some reason reading edited text from file adds new line. Strip it.
	plaintext = bytes.TrimRight(plaintext, "\n")

	// Encrypt edited password
	ciphertext, err = Encrypt(k, plaintext)
	if err != nil {
		log.Fatal("unable to re-encrypt password file after editing | ", err)
	}

	// Write encrypted password back to file
	err = ioutil.WriteFile(f, ciphertext, 0755)
	if err != nil {
		log.Fatal("unable to write re-encrypted password to file after editing | ", err)
	}
}
