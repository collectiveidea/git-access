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
			Name:  "authorize-command",
			Value: "git-access",
			Usage: "[Authorized Keys] Path to binary that will be inserted into the command option of the returned Authorized Keys.",
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
		var err error

		if c.Bool("authorized-keys") {
			err = authorizedKeysRequest(c)
		} else {
			err = gitRequest(c)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}

func authorizedKeysRequest(c *cli.Context) error {
	keysUrl := c.String("authorized-keys-url")

	if keysUrl == "" {
		return fmt.Errorf("The flag --authorized-keys-url is required when --authorized-keys/-A is used. See --help for more info.")
	}

	return RequestAuthorizedKeys(c.String("authorize-command"), keysUrl)
}

func gitRequest(c *cli.Context) error {
	permissionCheckUrl := c.String("permission-check-url")
	if permissionCheckUrl == "" {
		return fmt.Errorf("Missing required parameter --permission-check-url. See --help for more info.")
	}

	userId := c.String("user")
	if userId == "" {
		return fmt.Errorf("Missing required parameter --user. See --help for more info.")
	}

	sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND")
	if sshCommand == "" {
		return fmt.Errorf("No ssh command found")
	}

	return RequestGitAccess(sshCommand, userId, permissionCheckUrl)
}
