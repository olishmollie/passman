package passman

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

// Edit decrypts and opens password file in editor defined by $VISUAL
func Edit(p string) {
	f := path.Join(Root, p)
	if !pswdExists(f) {
		FatalError(nil, "could not find password for "+p)
	}

	// Read encrypted password from file
	ciphertext, err := ioutil.ReadFile(f)
	if err != nil {
		FatalError(err, "could not read pswd for "+p)
	}

	// Decrypt password
	plaintext, err := decrypt(getEncryptionKey(), ciphertext)
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
	ciphertext, err = encrypt(getEncryptionKey(), plaintext)
	if err != nil {
		FatalError(err, "could not encrypt pswd for "+p+" after editing")
	}

	// Write encrypted password back to file
	err = ioutil.WriteFile(f, ciphertext, 0755)
	if err != nil {
		FatalError(err, "could not write encrypted pswd to file after editing")
	}
}
