package main

import (
	"fmt"
	"os"
	"path"

	"github.com/atotto/clipboard"

	"github.com/docopt/docopt-go"
	"github.com/olishmollie/passman/passman"
)

var version = "v0.5"
var usage = `Usage:
	passman
	passman [-c] <prefix>
	passman add <prefix> <password>
	passman delete <prefix>
	passman dump [-o <outfile>]
	passman edit <prefix>
	passman generate [-cn] [-l int]
	passman import <infile>
	passman init
	passman nuke [-f]
	passman -h | --help
	passman -v | --version

Options:
	-c, --copy                Copy to clipboard. 	  
	-f, --force               Nuke w/o confirmation.
	-h, --help                Show this screen.
	-l, --length=<int>        Specify length of generated password.
	-n, --nosym               Generate password w/ no symbols.
	-o, --out=<outfile>       Specify file to be written to [default: pswds~].
	-v, --version             Show version.
`

func main() {

	args, err := docopt.ParseDoc(usage)
	if err != nil {
		passman.FatalError(err, "unable to parse command line options")
	}

	root := passman.GetRootDir()
	keyfile := path.Join(root, ".key")

	prefix, _ := args.String("<prefix>")
	password, _ := args.String("<password>")
	infile, _ := args.String("<infile>")
	outfile, _ := args.String("--out")
	copy, _ := args.Bool("--copy")
	nosym, _ := args.Bool("--nosym")
	force, _ := args.Bool("--force")
	len, _ := args.Int("--length")

	switch {
	case args["--version"]:
		fmt.Println(version)
	case args["--help"]:
		fmt.Println(usage)
	case args["init"]:
		passman.Init(root, keyfile)
	case args["add"]:
		check(root, keyfile)
		passman.Add(root, keyfile, prefix, password)
	case args["delete"]:
		check(root, keyfile)
		passman.Remove(root, prefix)
	case args["edit"]:
		check(root, keyfile)
		passman.Edit(root, keyfile, prefix)
	case args["generate"]:
		pswd := passman.Generate(len, nosym)
		copyOrPrint(pswd, copy)
	case args["dump"]:
		check(root, keyfile)
		passman.Dump(root, keyfile, outfile)
	case args["import"]:
		checkStore(root)
		passman.Import(root, keyfile, infile)
	case args["nuke"]:
		check(root, keyfile)
		passman.Nuke(root, force)
	case args["<prefix>"] != nil:
		check(root, keyfile)
		p, isDir := passman.Find(root, keyfile, prefix)
		if isDir {
			passman.Print(p, 0)
		} else {
			copyOrPrint(p, copy)
		}
	default:
		check(root, keyfile)
		passman.Print(root, 0)
	}
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
