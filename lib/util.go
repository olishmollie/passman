package lib

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

// GetRootDir returns passman's root directory
func GetRootDir() string {
	h := getHomeDir()
	return path.Join(h, ".passman")
}

func getHomeDir() string {
	u := getUser()
	return u.HomeDir
}

func getUser() *user.User {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// PswdExists returns true if password file exists, false otherwise
func PswdExists(p string) bool {
	dir, file := SplitDir(p)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("error finding password | ", err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if f.Name() == file {
			return true
		}
	}
	return false
}

// DirExists returns true if directory exists, false otherwise
func DirExists(d string) bool {
	_, err := os.Stat(d)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatal("error finding directory | ", err)
	}
	return true
}

// Getln reads a newline delimted string from stdin
func Getln() []byte {
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	out := bytes.TrimRight(in, "\n")
	return out
}

// SplitDir splits a directory into its directory and file components
func SplitDir(p string) (dir, file string) {
	sp := strings.SplitAfter(p, "/")
	dir = strings.Join(sp[:len(sp)-1], "")
	file = sp[len(sp)-1]
	return
}
