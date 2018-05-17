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
	root := GetRootDir()
	fname := path.Join(root, dir)
	if !PswdExists(fname) {
		FatalError(nil, "cannot find pswd for "+fname)
	}
	ct, err := ioutil.ReadFile(fname)
	if err != nil {
		FatalError(err, "could not read pswd for "+fname)
	}
	k := getKey()
	pswd, err := Decrypt(k, ct)
	if err != nil {
		FatalError(err, "could not decrypt pswd for "+fname)
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
			FatalError(err, "could not create password store")
		}
	}

	if PswdExists(path.Join(root, p)) {
		FatalError(nil, "that password already exists. Try `passman edit`.")
	}

	f, err := os.Create(path.Join(root, dir, file))
	if err != nil {
		FatalError(err, "could not create password")
	}

	k := getKey()
	ct, err := Encrypt(k, []byte(data))
	if err != nil {
		FatalError(err, "could not encrypt password for "+dir)
	}

	_, err = f.Write(ct)
	if err != nil {
		FatalError(err, "could not write pswd for "+dir)
	}

}

// Remove removes given password from storage
func Remove(p string) {
	root := GetRootDir()
	dir := path.Join(root, p)
	if DirExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			FatalError(err, "could not remove pswd for "+dir)
		}
	}
}

// Edit decrypts and opens password file in editor defined by $VISUAL
func Edit(p string) {
	root := GetRootDir()
	f := path.Join(root, p)
	if !PswdExists(f) {
		FatalError(nil, "could not find password for "+p)
	}

	k := getKey()

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
func Dump(dir, outfile string) {
	if _, err := os.Stat(outfile); err == nil {
		err = os.Remove(outfile)
		if err != nil {
			FatalError(err, "could not remove existing outfile")
		}
	}
	dumpRec(dir, outfile)
}

func dumpRec(dir, outfile string) {
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
			dumpRec(p, outfile)
		} else {
			p := path.Join(dir, n)
			k := getKey()
			c, err := ioutil.ReadFile(p)
			if err != nil {
				FatalError(err, "could not read pswd for "+p)
			}
			t, err := Decrypt(k, c)
			if err != nil {
				FatalError(err, "could not decrypt pswd for "+dir)
			}
			f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				FatalError(err, "could not open dump file "+outfile)
			}
			root := GetRootDir()
			w := strings.TrimPrefix(p, root+"/") + " " + string(t) + "\n"
			_, err = f.WriteString(w)
			if err != nil {
				FatalError(err, "could not write to dump file")
			}
			f.Close()
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

// Lock dumps the password store into an encrypted file and removes all passwords and .fpubkey
// CAUTION: if you forget the password you used to generate the encryption key, you will not
// be able to unencrypt your passwords.
func Lock() {
	yes := Getln("CAUTION: if you forget the password you used to generate your encryption key,\nyou will not be able to unencrypt your passwords.\nDo you wish to continue? (y/N) ")
	if string(yes) == "y" || string(yes) == "Y" {
		fmt.Println("Locking passman...")
		k := getKey()
		root := GetRootDir()
		lockfile := path.Join(root, "passman.lock")
		Dump(root, lockfile)
		EncryptFile(k, lockfile)
		RemoveContents(root, "passman.lock")
	} else {
		fmt.Println("Lock operation aborted.")
	}
}

// Unlock takes the passman.lock file, unencrypts it and imports all passwords into the password store
func Unlock() {
	pswd := getUserKey()
	writeUserKey(pswd)
	root := GetRootDir()
	lockfile := path.Join(root, "passman.lock")
	k := getKey()
	DecryptFile(k, lockfile)
	Import(lockfile)
	os.Remove(lockfile)
}
