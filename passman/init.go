package passman

import "os"

// Init initializes passman by creating a storage directory and generating a cipher key
func Init(root, keyfile string) {
	if !DirExists(root) {
		err := os.Mkdir(root, 0755)
		if err != nil {
			FatalError(err, "could not create pswd store")
		}
		Log("created password store at ~/.passman")
	}

	if _, err := os.Stat(keyfile); err == nil {
		FatalError(nil, "encryption key detected. remove `~/.passman/.key` before reinitializing")
	}

	key := generateEncryptionKey()
	Log("writing encryption key to .passman/.key...")
	writeEncryptionKey(keyfile, key)
	Log("passman initialized successfully")
}

func writeEncryptionKey(keyfile string, key []byte) {
	f, err := os.Create(keyfile)
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
	return []byte(Generate(32, false, false))
}
