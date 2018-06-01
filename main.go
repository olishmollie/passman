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
		checkNumArgs(2, args, cmd)
		passman.Add(root, keyfile, args[0], args[1])
	case "rm":
		checkLock(lockfile)
		checkStore(root)
		checkNumArgs(1, args, cmd)
		passman.Remove(root, args[0])
	case "edit":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(1, args, cmd)
		passman.Edit(root, keyfile, args[0])
	case "generate":
		checkNumArgs(0, args, cmd)
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
		checkNumArgs(0, args, cmd)
		passman.Dump(root, root, keyfile, os.Stdout)
	case "import":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(1, args, cmd)
		passman.Import(root, keyfile, args[0])
	case "lock":
		checkLock(lockfile)
		checkStore(root)
		checkFPubKey(keyfile)
		checkNumArgs(0, args, cmd)
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
		checkNumArgs(0, args, cmd)
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
		passman.FatalError(nil, "no password store. try `passman init`")
	}
}

func checkFPubKey(keyfile string) {
	if _, err := os.Stat(keyfile); err != nil {
		if os.IsNotExist(err) {
			passman.FatalError(nil, "no encryption key. try `passman init`")
		} else {
			passman.FatalError(err, "could not check status of password store")
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

func checkNumArgs(num int, args []string, cmd string) {
	if len(args) != num {
		fmt.Fprintf(os.Stderr, "incorrect number of arguments for %s. try `passman -h`\n", cmd)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(usage)
}

var usage = `usage: passman [-c] [-n] [-l int] [command]

commands: dump edit generate import init rm touch

passman - displays all passwords in store

passman [opts...] <prefix> - prints unencrypted password to stdout
    -c, --copy
        copies password to clipboard

add <prefix> <password> - add <password> to <prefix> directory

dump - prints unencrypted passwords to stdout

edit <prefix> - edit password(s) in <prefix> directory with default editor

generate [opts...] - generates a random password
    -c, --copy
        copies unencrypted password to clipboard
    -l, --len int 
        specifies length of generated password
    -n, --nosym 
        generate a password with no symbols

import <infile> - imports passwords from infile.
    NOTE: infile must be in the following format:
        website/username secret_password
        Category/anothersite/username another_password
        etc.

init - create password store if it doesn't exist, generate encryption key

lock - asks for a password, generates a symmetric cipher, encrypts and dumps store into ~/.passman/.passman.lock.
	NOTE: passman commands will be unavailable until running 'passman unlock'. this is one way to migrate your passwords.
	CAUTION: if you forget the password your provide, you will not be able to recover your passwords

rm <prefix> - remove password(s) in <prefix> directory from password store

unlock - undoes lock operation
`
