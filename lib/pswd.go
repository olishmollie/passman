package lib

import (
	"strings"
)

// Pswd encapsulates password data to be written to a file
type Pswd struct {
	path string
	data []byte
}

// NewPswd returns a new password with the given data
func NewPswd(path string, data []byte) Pswd {
	return Pswd{path, data}
}

// GetCg returns password category
func (p Pswd) GetCg() string {
	return strings.Split(p.path, "/")[0]
}

// GetAcct returns password acct
func (p Pswd) GetAcct() string {
	return strings.Split(p.path, "/")[1]
}

// GetData returns password data
func (p Pswd) GetData() []byte {
	return p.data
}
