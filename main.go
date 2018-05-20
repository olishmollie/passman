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

	if len(args) == 0 {
		checkLock()
		checkStore()
		passman.Print(passman.Root, 0)
		os.Exit(0)
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "init":
		checkLock()
		passman.Init()
	case "add":
		checkLock()
		checkStore()
		checkFPubKey()
		checkNumArgs(2, args)
		passman.Add(args[0], args[1])
	case "rm":
		checkLock()
		checkStore()
		checkNumArgs(1, args)
		passman.Remove(args[0])
	case "edit":
		checkLock()
		checkStore()
		checkFPubKey()
		checkNumArgs(1, args)
		passman.Edit(args[0])
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
		checkLock()
		checkStore()
		checkFPubKey()
		checkNumArgs(0, args)
		passman.Dump(passman.Root, os.Stdout)
	case "import":
		checkLock()
		checkStore()
		checkFPubKey()
		checkNumArgs(1, args)
		passman.Import(args[0])
	case "lock":
		checkLock()
		checkStore()
		checkFPubKey()
		checkNumArgs(0, args)
		passman.Lock()
	case "unlock":
		if locked() {
			passman.Unlock()
		} else {
			fmt.Println("passman is not locked")
		}
	default:
		checkLock()
		checkFPubKey()
		checkStore()
		checkNumArgs(0, args)
		p := passman.Find(cmd)
		if *copyPtr {
			clipboard.WriteAll(p)
		} else {
			fmt.Println(p)
		}
	}

}

func checkStore() {
	if !passman.DirExists(passman.Root) {
		passman.FatalError(nil, "no pswd store. try `passman init`")
	}
}

func checkFPubKey() {
	fpubkey := path.Join(passman.Root, ".key")
	if _, err := os.Stat(fpubkey); err != nil {
		if os.IsNotExist(err) {
			passman.FatalError(nil, "no encryption key. try `passman init`")
		} else {
			passman.FatalError(err, "could not check status of pswd store")
		}
	}
}

func checkLock() {
	if locked() {
		fmt.Println("passman is locked. try `passman unlock`")
		os.Exit(0)
	}
}

func locked() bool {
	if _, err := os.Stat(passman.Lockfile); err == nil {
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
