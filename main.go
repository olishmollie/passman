package main

import (
	"fmt"
	"os"

	"github.com/olishmollie/passman/lib"
)

func main() {

	args := os.Args[1:]

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
		p := lib.NewPswd(args[0], []byte(args[1]))
		lib.Add(p)
	case "rm":
		// TODO: Remove password
	default:
		checkStore()
		checkNumArgs(1, args)
		lib.Find(args[0])
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
