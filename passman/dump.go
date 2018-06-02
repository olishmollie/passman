package passman

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Dump re-encrypts passwords using a user provided password and writes them to outfile.
func Dump(root, keyfile, outfile string) {

	pswd := getUserPswd()
	k := hashPswd(pswd)

	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		FatalError(err, "could not open outfile")
	}
	defer f.Close()

	writePswdStore(root, root, keyfile, f)

	encryptFile(k[:], outfile)
}

func getUserPswd() []byte {
	var pswd []byte
	// TODO - limit number of times passphrase can be entered
	for {
		pswd = getln("Enter your passphrase: ")
		pswd2 := getln("Confirm passphrase: ")
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

func writePswdStore(start, root, keyfile string, out *os.File) {

	files, err := ioutil.ReadDir(start)
	if err != nil {
		FatalError(err, "could not read files from password store")
	}

	for _, file := range files {

		n := file.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}

		if file.IsDir() {

			p := path.Join(start, n)
			writePswdStore(p, root, keyfile, out)

		} else {

			p := path.Join(start, n)
			c, err := ioutil.ReadFile(p)
			if err != nil {
				FatalError(err, "could not read pswd for "+p)
			}

			t, err := decrypt(getEncryptionKey(keyfile), c)
			w := strings.TrimPrefix(p, root+"/") + " " + string(t) + "\n"
			if err != nil {
				FatalError(err, "could not decrypt pswd for "+w)
			}

			_, err = out.WriteString(w)
			if err != nil {
				FatalError(err, "could not write passwords to tmp file")
			}
		}
	}
}
