package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const lbar = "|\u2014\u2014 "

// Print prints a tree of the password store
func Print(dirName string, offset int) {
	files, err := ioutil.ReadDir(dirName)
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
			p := path.Join(dirName, n)
			Print(p, offset+2)
		}
	}
}

// Find finds, decrypts, and prints a password to the console
func Find(dir string) {
	root := GetRootDir()
	fname := path.Join(root, dir)
	if !DirExists(fname) {
		fmt.Println("error: that password doesn't exist. did you remember a category prefix?")
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
	fmt.Println(string(pswd))
}

// Add inserts a password into storage
func Add(p Pswd) {

	root := GetRootDir()
	cgdir := path.Join(root, p.GetCg())
	if !DirExists(cgdir) {
		err := os.Mkdir(cgdir, 0755)
		if err != nil {
			log.Fatal("error creating cgDir | ", err)
		}
	}

	acct := path.Join(cgdir, p.GetAcct())
	f, err := os.Create(acct)
	if err != nil {
		log.Fatal("error making password file | ", err)
	}

	k := GetKey()
	ct, err := Encrypt(k, []byte(p.data))
	if err != nil {
		log.Fatal("error encrypting password data | ", err)
	}

	_, err = f.Write(ct)
	if err != nil {
		log.Fatal("error writing to file | ", err)
	}

}
