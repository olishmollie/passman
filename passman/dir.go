package passman

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// FileExists returns whether file specifed by name exists.
func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		FatalError(err, "could not find file "+name)
	}
	return true
}

// GetRootDir returns ~/.passman
func GetRootDir(rootName string) string {
	h := getHomeDir()
	return path.Join(h, rootName)
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

func removeContentsOf(dir string, except ...string) {
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
