package util

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(name string, args ...string) {
	c := exec.Command(name, args...)
	out, err := c.Output()
	if err != nil {
		fmt.Printf("Failed to run %s %v: %+v", name, args, err)
		os.Exit(1)
	}
	if len(out) > 0 {
		fmt.Println(string(out))
	}
}
