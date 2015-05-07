package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	var app = cli.NewApp()
	app.Name = "git-access"
	app.Usage = "Protect access to Git repositories over SSH"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "authorized-keys,A",
			Usage: "Toggle Authorized Keys mode",
		},
		cli.StringFlag{
			Name:  "authorized-keys-url",
			Value: "",
			Usage: "HTTP(S) Endpoint for querying valid public SSH keys. Only valid when using -A.",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("authorized-keys") {
			readAuthorizedKeys(c.String("authorized-keys-url"))
		} else {
			processGitRequest(c.Args())
		}
	}

	app.Run(os.Args)
}

func readAuthorizedKeys(keysUrl string) {
	if keysUrl == "" {
		fmt.Println("Error: --authorized-keys-url is required when --authorized-keys is used")
		os.Exit(1)
	}

	response, err := http.Get(keysUrl)

	if err != nil {
		fmt.Println("Net Error:", err)
		return
	}
	defer response.Body.Close()

	keys, _ := ioutil.ReadAll(response.Body)

	for _, key := range strings.Split(string(keys), "\n") {
		fmt.Println(
			"command=\"git-access\",no-user-rc,no-X11-forwarding,no-agent-forwarding,no-pty",
			key,
		)
	}
}

func processGitRequest(args []string) {
	action := args[0]
	fullActionPath, err := exec.LookPath(action)
	if err != nil {
		fmt.Println("Unable to find the binary", action)
		os.Exit(1)
	}

	if action == "git-receive-pack" ||
		action == "git-upload-pack" ||
		action == "git-upload-archive" {

		err = syscall.Exec(fullActionPath, args, []string{})
		fmt.Println("Woah, failed to execute the action", action, err)
	} else {
		os.Exit(1)
	}
}
