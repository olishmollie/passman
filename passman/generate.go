package passman

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"1234567890" +
	"!@#$%^&*()[]{}<>-_=+,.?/'\":;"

const charsetNoSym = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"1234567890"

// Generate takes an options hash and returns a randomly generated password
// If l == 0, a random length will be provided between 8 and 20
func Generate(l int, noSym bool) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	if l == 0 {
		l = seededRand.Intn(13) + 8
	}
	b := make([]byte, l)
	for i := range b {
		if noSym {
			b[i] = charsetNoSym[seededRand.Intn(len(charsetNoSym))]
		} else {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
	}
	return string(b)
}
