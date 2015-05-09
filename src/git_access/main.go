package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	var app = cli.NewApp()
	app.Name = "git-access"
	app.Usage = "Protect access to Git repositories over SSH"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "authorized-keys,A",
			Usage: "Toggle Authorized Keys mode. If not set will be in Git Access mode.",
		},
		cli.StringFlag{
			Name:  "authorized-keys-url",
			Usage: "[Authorized Keys] HTTP(S) Endpoint for querying valid public SSH keys. Only valid when using -A.",
		},
		cli.StringFlag{
			Name:  "user,U",
			Usage: "[Git Access] Unique User identifier for git access permissions check.",
		},
		cli.StringFlag{
			Name:  "permission-check-url",
			Usage: "[Git Access] HTTP(S) Endpoint for querying repository permissions.",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("authorized-keys") {
			authorizedKeysRequest(c)
		} else {
			gitRequest(c)
		}
	}

	app.Run(os.Args)
}

func authorizedKeysRequest(c *cli.Context) {
	keysUrl := c.String("authorized-keys-url")

	if keysUrl == "" {
		fmt.Println("Error: --authorized-keys-url is required when --authorized-keys/-A is used")
		os.Exit(1)
	}

	err := RequestAuthorizedKeys(keysUrl)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func gitRequest(c *cli.Context) {
	permissionsUrl := c.String("permission-check-url")
	if permissionsUrl == "" {
		fmt.Println("Missing required parameter --permission-check-url")
		os.Exit(1)
	}

	userId := c.String("user")
	if userId == "" {
		fmt.Println("Missing required parameter --user")
		os.Exit(1)
	}

	sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND")
	if sshCommand == "" {
		fmt.Println("No ssh command found")
		os.Exit(1)
	}

	GitRequest(userId, permissionsUrl)
}
