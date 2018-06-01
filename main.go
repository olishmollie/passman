package main

import (
	"fmt"
	"os"
	"path"

	"github.com/atotto/clipboard"

	"github.com/docopt/docopt-go"
	"github.com/olishmollie/passman/passman"
)

var version = "v0.4"
var usage = `Usage:
	passman
	passman [-c] <prefix>
	passman add <prefix> <password>
	passman delete <prefix>
	passman dump
	passman edit <prefix>
	passman generate [-cn] [-l int]
	passman import <infile>
	passman init
	passman lock
	passman unlock
	passman -h | --help
	passman -v | --version

Options:
	-h, --help               Show this screen.
	-v, --version            Show version.
	-c, --copy               Copy to clipboard. 	  
	-n, --nosym              Generate password w/ no symbols.
	-l int, --length=int     Specify length of generated password.
`

func main() {

	args, err := docopt.ParseDoc(usage)
	if err != nil {
		passman.FatalError(err, "unable to parse command line options")
	}

	root := passman.GetRootDir()
	keyfile := path.Join(root, ".key")
	lockfile := path.Join(root, ".passman.lock")

	prefix, _ := args.String("<prefix>")
	password, _ := args.String("<password>")
	infile, _ := args.String("<infile>")
	copy, _ := args.Bool("--copy")
	nosym, _ := args.Bool("--nosym")
	len, _ := args.Int("--length")

	locked := isLocked(lockfile)

	if !locked {
		check(root, keyfile)
		switch {
		case args["--version"]:
			fmt.Println(version)
		case args["--help"]:
			fmt.Println(usage)
		case args["init"]:
			passman.Init(root, keyfile)
		case args["add"]:
			passman.Add(root, keyfile, prefix, password)
		case args["delete"]:
			passman.Remove(root, prefix)
		case args["edit"]:
			passman.Edit(root, keyfile, prefix)
		case args["generate"]:
			pswd := passman.Generate(len, nosym)
			copyOrPrint(pswd, copy)
		case args["dump"]:
			passman.Dump(root, root, keyfile, os.Stdout)
		case args["import"]:
			passman.Import(root, keyfile, infile)
		case args["lock"]:
			passman.Lock(root, keyfile, lockfile)
		case args["unlock"]:
			passman.FatalError(nil, "passman is not locked")
		case args["<prefix>"] != nil:
			p, isDir := passman.Find(root, keyfile, prefix)
			if isDir {
				passman.Print(p, 0)
			} else {
				copyOrPrint(p, copy)
			}
		default:
			passman.Print(root, 0)
		}
	} else {
		switch {
		case args["unlock"]:
			passman.Unlock(root, keyfile, lockfile)
		default:
			passman.FatalError(nil, "passman is locked. try `passman unlock`")
		}
	}
}

func isLocked(lockfile string) bool {
	if _, err := os.Stat(lockfile); err == nil {
		return true
	}
	return false
}

func check(root, keyfile string) {
	checkStore(root)
	checkKey(keyfile)
}

func checkStore(root string) {
	if !passman.DirExists(root) {
		passman.FatalError(nil, "no password store. try `passman init`")
	}
}

func checkKey(keyfile string) {
	if _, err := os.Stat(keyfile); err != nil {
		if os.IsNotExist(err) {
			passman.FatalError(nil, "no encryption key. try `passman init`")
		} else {
			passman.FatalError(err, "could not check status of password store")
		}
	}
}

func copyOrPrint(s string, copy bool) {
	if copy {
		clipboard.WriteAll(s)
	} else {
		fmt.Println(s)
	}
}
