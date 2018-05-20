package passman

import (
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

// DirExists returns whether given directory exists.
func DirExists(d string) bool {
	_, err := os.Stat(d)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		FatalError(err, "could not find dir "+d)
	}
	return true
}

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

func pswdExists(p string) bool {
	dir, file := splitDir(p)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not read files in pswd dir "+p)
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

func splitDir(p string) (dir, file string) {
	sp := strings.SplitAfter(p, "/")
	dir = strings.Join(sp[:len(sp)-1], "")
	file = sp[len(sp)-1]
	return
}

func removeContents(dir string, except ...string) {
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
