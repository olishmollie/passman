package lib

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/user"
	"path"
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
