package passman

import (
	"os"
	"os/exec"
)

// Print displays all the password prefixes in dir
func Print(dir string, offset int) {
	cmd := exec.Command("tree", "--noreport", dir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		FatalError(err, "could not open editor assigned to $VISUAL")
	}
}
