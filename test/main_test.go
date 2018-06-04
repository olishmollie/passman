package main

import (
	"io"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/olishmollie/passman/passman"
)

type test struct {
	name string
	cmd  *exec.Cmd
}

func (t *test) init(in, out io.ReadWriter, args ...string) {
	t.name = args[0]
	t.cmd = exec.Command("./passman", args...)
	t.cmd.Stderr = os.Stderr
	t.cmd.Stdin = in
	t.cmd.Stdout = out
}

func (t *test) run() error {
	return t.cmd.Run()
}

func TestMain(m *testing.M) {
	build()
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func build() {
	make := exec.Command("make", "-C", "../", "build_test")
	err := make.Run()
	if err != nil {
		passman.FatalError(err, "unable to build test binary")
	}
}

var root = passman.GetRootDir(".passman_test")
var keyfile = path.Join(root, ".key")
var key = passman.GenerateEncryptionKey()

func setup() {

	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(root, 0777)
			if err != nil {
				passman.FatalError(err, "unable to create mock password store")
			}
		}
	}

	f, err := os.OpenFile(keyfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		passman.FatalError(err, "unable to open mock keyfile")
	}
	defer f.Close()

	f.Write(key)
}

func teardown() {
	err := os.RemoveAll(root)
	if err != nil {
		passman.FatalError(err, "unable to remove mock store")
	}
	err = os.RemoveAll("passman")
	if err != nil {
		passman.FatalError(err, "unable to remove test bin")
	}
}
