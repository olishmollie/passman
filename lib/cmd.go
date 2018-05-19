package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Init initializes passman by creating a storage directory and generating a cipher key
func Init() {
	if !DirExists(Root) {
		err := os.Mkdir(Root, 0755)
		if err != nil {
			FatalError(err, "could not create pswd store")
		}
		Log("created password store at ~/.passman")
	}

	if _, err := os.Stat(Keyfile); err == nil {
		FatalError(nil, "encryption key detected. remove `~/.passman/.key` before reinitializing")
	}

	key := generateEncryptionKey()
	Log("writing encryption key to .passman/.key...")
	writeEncryptionKey(key)
	Log("passman initialized successfully")
}

func writeEncryptionKey(key []byte) {
	f, err := os.Create(Keyfile)
	if err != nil {
		FatalError(err, "could not create key file")
	}
	_, err = f.Write(key[:])
	if err != nil {
		FatalError(err, "could not write key to key file")
	}
	err = f.Close()
	if err != nil {
		FatalError(err, "could not close key file")
	}
}

func generateEncryptionKey() []byte {
	return []byte(Generate(32, false))
}

// Print prints a tree of the password store
func Print(dir string, offset int) {
	const lbar = "|\u2014\u2014 "
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not read files in password store")
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
	fname := path.Join(Root, dir)
	if !PswdExists(fname) {
		FatalError(nil, "cannot find pswd for "+dir)
	}
	ct, err := ioutil.ReadFile(fname)
	if err != nil {
		FatalError(err, "could not read pswd for "+dir)
	}
	k := getEncryptionKey()
	pswd, err := Decrypt(k, ct)
	if err != nil {
		FatalError(err, "bad encryption key for "+dir)
	}
	return string(pswd)
}

// Add inserts a password into storage
func Add(p, data string) {

	dir, file := SplitDir(p)

	newdir := path.Join(Root, dir)
	if !DirExists(newdir) {
		err := os.MkdirAll(newdir, 0755)
		if err != nil {
			FatalError(err, "could not create password store")
		}
	}

	f, err := os.Create(path.Join(Root, dir, file))
	if err != nil {
		FatalError(err, "could not create password")
	}

	k := getEncryptionKey()
	ct, err := Encrypt(k, []byte(data))
	if err != nil {
		FatalError(err, "could not encrypt password for "+dir)
	}

	_, err = f.Write(ct)
	if err != nil {
		FatalError(err, "could not write pswd for "+dir)
	}

	err = f.Close()
	if err != nil {
		FatalError(err, "could not close pswd for "+dir)
	}

}

// Remove removes given password from storage
func Remove(p string) {
	dir := path.Join(Root, p)
	if DirExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			FatalError(err, "could not remove pswd for "+dir)
		}
	}
}

// Edit decrypts and opens password file in editor defined by $VISUAL
func Edit(p string) {
	f := path.Join(Root, p)
	if !PswdExists(f) {
		FatalError(nil, "could not find password for "+p)
	}

	k := getEncryptionKey()

	// Read encrypted password from file
	ciphertext, err := ioutil.ReadFile(f)
	if err != nil {
		FatalError(err, "could not read pswd for "+p)
	}

	// Decrypt password
	plaintext, err := Decrypt(k, ciphertext)
	if err != nil {
		FatalError(err, "could not decode pswd for "+p)
	}

	// Write plaintext to password file
	err = ioutil.WriteFile(f, plaintext, 0775)
	if err != nil {
		FatalError(err, "could not write pswd for "+p)
	}

	// Open in editor
	cmd := exec.Command(os.ExpandEnv("$VISUAL"), f)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		FatalError(err, "could not open editor assigned to $VISUAL")
	}

	// Read edited text from file
	plaintext, err = ioutil.ReadFile(f)
	if err != nil {
		FatalError(err, "could not read file after editing")
	}

	// For some reason reading edited text from file adds new line. Strip it.
	plaintext = bytes.TrimRight(plaintext, "\n")

	// Encrypt edited password
	ciphertext, err = Encrypt(k, plaintext)
	if err != nil {
		FatalError(err, "could not encrypt pswd for "+p+" after editing")
	}

	// Write encrypted password back to file
	err = ioutil.WriteFile(f, ciphertext, 0755)
	if err != nil {
		FatalError(err, "could not write encrypted pswd to file after editing")
	}
}

// Dump writes unecrypted passwords to outfile.
func Dump(dir string, outfile *os.File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		FatalError(err, "could not read files from password store")
	}
	for _, f := range files {
		n := f.Name()
		if strings.HasPrefix(n, ".") {
			continue
		}
		if f.IsDir() {
			p := path.Join(dir, n)
			Dump(p, outfile)
		} else {
			p := path.Join(dir, n)
			k := getEncryptionKey()
			c, err := ioutil.ReadFile(p)
			if err != nil {
				FatalError(err, "could not read pswd for "+p)
			}
			t, err := Decrypt(k, c)
			w := strings.TrimPrefix(p, Root+"/") + " " + string(t) + "\n"
			if err != nil {
				FatalError(err, "could not decrypt pswd for "+w)
			}
			_, err = outfile.WriteString(w)
			if err != nil {
				FatalError(err, "could not write passwords to dump file")
			}
		}
	}
}

// Import reads pswds from a file and adds them to pswd store.
// File must be newline delimited directory and pswds (like what Dump outputs)
// e.g. Category/Website/username pswd
//      Category/NextWebsite/username pswd
func Import(infile string) {
	f, err := os.Open(infile)
	if err != nil {
		FatalError(err, "could not open import file")
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		d := strings.Split(t, " ")
		Add(d[0], d[1])
	}
}

// Lock dumps the password store into an encrypted file and removes all passwords and .key
// CAUTION: if you forget the password you used to generate the encryption key, you will not
// be able to unencrypt your passwords.
func Lock() {
	Log("CAUTION - if you forget the password you use to lock passman, you will not be able to unencrypt your passwords.")
	yes := Getln("Do you wish to continue? (y/N) ")
	if string(yes) == "y" || string(yes) == "Y" {
		pswd := getUserPswd()
		k := hashPswd(pswd)
		lockfile := path.Join(Root, ".passman.lock")
		f, err := os.OpenFile(lockfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			FatalError(err, "could not open lockfile")
		}
		Dump(Root, f)
		EncryptFile(k[:], lockfile)
		Log("deleting password store...")
		RemoveContents(Root, ".passman.lock")
		Log("passman locked")
		err = f.Close()
		if err != nil {
			FatalError(err, "could not close lockfile")
		}
	} else {
		Log("lock operation aborted.")
	}
}

// Unlock unencrypts the .passman.lock file and imports all passwords into the password store
func Unlock() {
	pswd := getUserPswd()
	key := hashPswd(pswd)
	DecryptFile(key[:], Lockfile)
	newKey := generateEncryptionKey()
	writeEncryptionKey(newKey)
	Import(Lockfile)
	os.Remove(Lockfile)
}
