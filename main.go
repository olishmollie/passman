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
		fmt.Println("error: password store is empty. try `passman init`.")
		os.Exit(1)
	}
	return root
}

func checkNumArgs(num int, args []string) {
	if len(args) != num {
		fmt.Println("USAGE")
		os.Exit(1)
	}
}
