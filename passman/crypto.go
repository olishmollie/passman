package passman

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

func getEncryptionKey() []byte {
	key, err := ioutil.ReadFile(Keyfile)
	if err != nil {
		FatalError(err, "could not read encryption key from password store")
	}
	return key
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func encryptFile(key []byte, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		FatalError(err, "could not read data from file to be encrypted")
	}
	cipher, err := encrypt(key, data)
	if err != nil {
		FatalError(nil, "could not encrypt file data")
	}
	err = ioutil.WriteFile(filename, cipher, 0666)
	if err != nil {
		FatalError(err, "could not write encrypted data to file")
	}
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func decryptFile(key []byte, filename string) {
	cipher, err := ioutil.ReadFile(filename)
	if err != nil {
		FatalError(err, "could not read data from file to be decrypted")
	}
	data, err := decrypt(key, cipher)
	if err != nil {
		FatalError(nil, "could not decrypt file data. incorrect key")
	}
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		FatalError(err, "could not write encrypted data to file")
	}
}
