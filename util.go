package passman

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// Root is equivalent to ~/.passman
var Root = getRootDir()

// Lockfile is the path to the lockfile
var Lockfile = path.Join(Root, ".passman.lock")

// Keyfile is the path to the keyfile
var Keyfile = path.Join(Root, ".key")

func getRootDir() string {
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
		FatalError(err, "could not find current user")
	}
	return u
}

// PswdExists returns true if password file exists, false otherwise
func PswdExists(p string) bool {
	dir, file := SplitDir(p)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not find pswd "+p)
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
		FatalError(err, "could not find dir "+d)
	}
	return true
}

// Getln reads a newline delimted string from stdin
func Getln(msg string) []byte {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadBytes('\n')
	if err != nil {
		FatalError(err, "could not read from stdin")
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

// FatalError writes an error to stderr and exits
func FatalError(err error, msg string) {
	if err == nil {
		os.Stderr.WriteString("fatal: " + msg + "\n")
		os.Exit(1)
	}
	os.Stderr.WriteString("fatal: " + msg + "\n" + err.Error() + "\n")
	os.Exit(1)
}

// Log writes a message stdout
func Log(msg string) {
	fmt.Printf("passman: %s\n", msg)
}

// RemoveContents removes a directory's contents, skipping any filenames passed as secondary arguments.
func RemoveContents(dir string, except ...string) {
	d, err := os.Open(dir)
	if err != nil {
		FatalError(err, "could not open pswd store")
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		FatalError(err, "could not read dirnames in pswd store")
	}
	for _, name := range names {
		skip := false
		for _, el := range except {
			if name == el {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			FatalError(err, "could not remove pswds from pswd store")
		}
	}
}

func getUserPswd() []byte {
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
	return bytes.TrimRight(pswd, "\n")
}

func hashPswd(p []byte) [32]byte {
	return sha256.Sum256(p)
}

func getEncryptionKey() []byte {
	d, err := ioutil.ReadFile(Keyfile)
	if err != nil {
		FatalError(err, "could not read encryption key from password store")
	}
	return d
}
