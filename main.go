package main

import (
	"fmt"
	"os"
	"path"

	"github.com/atotto/clipboard"
	"github.com/ogier/pflag"
	"github.com/olishmollie/passman/passman"
)

func main() {

	pflag.Usage = printUsage

	copyPtr := pflag.BoolP("copy", "c", false, "copy password to clipboard")
	noSymPtr := pflag.BoolP("nosym", "n", false, "generate a password with no symbols")
	lenPtr := pflag.IntP("len", "l", 0, "length of generated password, defaults to random int between 8 and 20")
	pflag.Parse()

	args := pflag.Args()

	root := passman.GetRootDir()
	keyfile := path.Join(root, ".key")
	lockfile := path.Join(root, ".passman.lock")

	if len(args) == 0 {
		checkLock(lockfile)
		checkStore(root)
		passman.Print(root, 0)
		os.Exit(0)
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "init":
		checkLock(lockfile)
		passman.Init(root, keyfile)
	case "add":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(2, args)
		passman.Add(root, keyfile, args[0], args[1])
	case "rm":
		checkLock(lockfile)
		checkStore(root)
		checkNumArgs(1, args)
		passman.Remove(root, args[0])
	case "edit":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(1, args)
		passman.Edit(root, keyfile, args[0])
	case "generate":
		checkNumArgs(0, args)
		var s string
		s = passman.Generate(*lenPtr, *noSymPtr)
		if *copyPtr {
			clipboard.WriteAll(s)
		} else {
			fmt.Println(s)
		}
	case "dump":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(0, args)
		passman.Dump(root, root, keyfile, os.Stdout)
	case "import":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(1, args)
		passman.Import(root, keyfile, args[0])
	case "lock":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(0, args)
		passman.Lock(root, keyfile, lockfile)
	case "unlock":
		if locked(lockfile) {
			passman.Unlock(root, keyfile, lockfile)
		} else {
			fmt.Println("passman is not locked")
		}
	default:
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(0, args)
		p := passman.Find(root, keyfile, cmd)
		if *copyPtr {
			clipboard.WriteAll(p)
		} else {
			fmt.Println(p)
		}
	}

}

func checkStore(root string) {
	if !passman.DirExists(root) {
		passman.FatalError(nil, "no pswd store. try `passman init`")
	}
}

func checkFPubKey(keyfile string) {
	if _, err := os.Stat(keyfile); err != nil {
		if os.IsNotExist(err) {
			passman.FatalError(nil, "no encryption key. try `passman init`")
		} else {
			passman.FatalError(err, "could not check status of pswd store")
		}
	}
}

func checkLock(lockfile string) {
	if locked(lockfile) {
		fmt.Println("passman is locked. try `passman unlock`")
		os.Exit(0)
	}
}

func locked(lockfile string) bool {
	if _, err := os.Stat(lockfile); err == nil {
		return true
	}
	return false
}

func checkNumArgs(num int, args []string) {
	if len(args) != num {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(usage)
}

var usage = `usage: passman [command] [args...] [opts...]

commands: dump edit generate import init rm touch

passman - prints a tree of pswds in store

passman [opts...] <pswd_file> - prints unencrypted pswd
    -c, --copy
        copy password to clipboard

add <pswd_file> - add <pswd_file> to pswd store

dump - prints unencrypted pswds to stdout

edit <pswd_file> - edit pswd in editor set to $VISUAL

generate [opts...] - generates a random pswd
    -c, --copy
        copies unencrypted pswd to clipboard
    -l, --len int 
        specifies length of generated pswd
    -n, --nosym 
        generate a password with no symbols

import <infile> - imports passwords from infile.
    NOTE: infile must be in the following format:
        website/username secret_password
        Category/anothersite/username another_password
        etc.

init - create pswd store if it doesn't exist, generate encryption key

lock - encrypts and dumps all passwords into one file

rm <pswd_file> - remove <pswd_file> from pswd store

unlock - undoes lock operation
`
