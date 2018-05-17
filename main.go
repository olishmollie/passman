package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/olishmollie/passman/lib"
)

func main() {

	copyPtr := flag.Bool("copy", false, "copy password to clipboard")
	noSymPtr := flag.Bool("nosym", false, "generate a password with no symbols")
	lenPtr := flag.Int("len", 0, "length of generated password, defaults to random int between 8 and 20")
	flag.Parse()
	var args []string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		args = append(args, arg)
	}
	args = args[1:]

	if len(args) == 0 {
		root := checkStore()
		lib.Print(root, 0)
		os.Exit(0)
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "init":
		lib.Init()
	case "touch":
		checkStore()
		checkNumArgs(2, args)
		lib.Add(args[0], args[1])
	case "rm":
		checkStore()
		checkNumArgs(1, args)
		lib.Remove(args[0])
	case "edit":
		checkStore()
		checkNumArgs(1, args)
		lib.Edit(args[0])
	case "generate":
		checkNumArgs(0, args)
		var s string
		s = lib.Generate(*lenPtr, *noSymPtr)
		if *copyPtr {
			clipboard.WriteAll(s)
		} else {
			fmt.Println(s)
		}
	case "dump":
		root := checkStore()
		checkNumArgs(1, args)
		lib.Dump(root, args[0])
	case "import":
		checkStore()
		checkNumArgs(1, args)
		lib.Import(args[0])
	default:
		checkStore()
		checkNumArgs(0, args)
		p := lib.Find(cmd)
		if *copyPtr {
			clipboard.WriteAll(p)
		} else {
			fmt.Println(p)
		}
	}

}

func checkStore() string {
	root := lib.GetRootDir()
	if !lib.DirExists(root) {
		lib.FatalError(nil, "no pswd store. Try `passman init`.")
	}
	return root
}

func checkNumArgs(num int, args []string) {
	if len(args) != num {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(
		`usage: passman [opts...] [command] [args...]
example: passman touch Category/Website/username pswd
commands: dump edit generate import init rm touch
	passman - prints a tree of pswds in store

passman [opts...] <pswd_file> - prints unencrypted pswd
    opts:
        -copy - copies unencrypted pswd to clipboard

dump <outfile> - prints unencrypted pswds to outfile

edit <pswd_file> - edit pswd in editor set to $VISUAL

generate [opts...] - generates a random pswd
	opts:
		-copy - copies unencrypted pswd to clipboard
		-len=int - specifies length of generated pswd
		-nosym - generate a password with no symbols

import <infile> - imports passwords from infile. infile must be newline delimited pswds

init - create password store if it doesn't exist, and generate encryption key

rm <pswd_file> - remove <pswd_file> from pswd store

touch <pswd_file> - add <pswd_file> to pswd store
`)
}
