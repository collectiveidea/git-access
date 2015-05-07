package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var app = cli.NewApp()
	app.Name = "git-access"
	app.Usage = "Protect access to Git repositories over SSH"

	app.Action = func(c *cli.Context) {
		action := c.Args()[0]

		if action == "git-receive-pack" ||
			action == "git-upload-pack" ||
			action == "git-upload-archive" {

			processGitRequest(action, c.Args())
		}

		os.Exit(1)
	}

	app.Run(os.Args)
}

func processGitRequest(action string, args []string) {
	fullActionPath, err := exec.LookPath(action)
	if err != nil {
		fmt.Println("Unable to find the binary", action)
		os.Exit(1)
	}

	err = syscall.Exec(fullActionPath, args, []string{})
	fmt.Println("Woah, failed to execute the action", action, err)
}
