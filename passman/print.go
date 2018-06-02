package passman

import (
	"fmt"
	"os"
	"os/exec"
)

// Print displays all the password prefixes in dir
func Print(dir string, offset int) {

	fmt.Println("Password Store")

	cmd1 := exec.Command("tree", "-l", "--noreport", dir)
	cmd2 := exec.Command("tail", "-n", "+2")

	var err error
	cmd2.Stdin, err = cmd1.StdoutPipe()
	if err != nil {
		FatalError(err, "unable to pipe tree into tail")
	}
	cmd2.Stdout = os.Stdout

	err = cmd2.Start()
	err = cmd1.Run()
	err = cmd2.Wait()
	if err != nil {
		FatalError(err, "unable to print tree")
	}
}
