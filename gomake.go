package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// go tool dist list
	var stderr, stdout bytes.Buffer
	env := append(os.Environ(),
		"GOOS="+os.Args[1],
		"GOARCH="+os.Args[2],
		// "LDFLAGS=\"-s -w\"",
		"CGO_ENABLE=0")
	// "CC=gcc",
	// "CGO_ENABLE=1")
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", os.Args[3])
	// cmd := exec.Command("go", "build", "test.go")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("%s\nStderr: %s", err, stderr.String())
		fmt.Println(err)
	}

	fmt.Println(stdout.String())
}
