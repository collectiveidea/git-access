package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	shellwords "github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"syscall"
)

// GitRequest takes the requested git command, ensures the user has
// permission to the given repository, and rewrites itself to let
// the git command with it's work (clone or push) to the right repository.
func GitRequest(c *cli.Context) {
	originalAction := os.Getenv("SSH_ORIGINAL_COMMAND")

	if originalAction == "" {
		fmt.Println("No original command specified. Exiting")
		os.Exit(1)
	}

	commandParts, _ := shellwords.Parse(originalAction)
	action := commandParts[0]

	fullActionPath, err := exec.LookPath(action)
	if err != nil {
		fmt.Println("Unable to find the binary", action)
		os.Exit(1)
	}

	if isValidAction(action) {
		err = syscall.Exec(fullActionPath, commandParts, []string{})
		fmt.Println("Failed to execute the command", commandParts, err)
	}

	os.Exit(1)
}

func isValidAction(action string) bool {
	return action == "git-receive-pack" ||
		action == "git-upload-pack" ||
		action == "git-upload-archive"
}
