package passman

import (
	"math/rand"
	"time"

	"github.com/atotto/clipboard"
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
func Generate(l int, noSym, copy bool) string {
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
	if copy {
		clipboard.WriteAll(string(b))
	}
	return string(b)
}
