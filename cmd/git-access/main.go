package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"

	"github.com/urfave/cli"
)

func main() {
	var app = cli.NewApp()
	app.Name = "git-access"
	app.Usage = "Protect access to Git repositories over SSH"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "syslog",
			Usage: "Enable logging to syslog",
		},
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
		if c.Bool("syslog") {
			writer, err := syslog.New(syslog.LOG_INFO, "git-access")
			if err != nil {
				fmt.Printf("Unable to enable syslog:\n")
				fmt.Println(err)

				os.Exit(1)
			}

			// Turn off go's own logging timestamps
			log.SetFlags(0)
			log.SetOutput(writer)
		}

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
		log.Fatalln("The flag --authorized-keys-url is required when --authorized-keys/-A is used. See --help for more info.")
	}

	RequestAuthorizedKeys(c.String("authorize-command"), keysUrl)
}

func gitRequest(c *cli.Context) {
	permissionCheckUrl := c.String("permission-check-url")
	if permissionCheckUrl == "" {
		log.Fatalln("Missing required parameter --permission-check-url. See --help for more info.")
	}

	userId := c.String("user")
	if userId == "" {
		log.Fatalln("Missing required parameter --user. See --help for more info.")
	}

	sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND")
	if sshCommand == "" {
		log.Fatalln("No original ssh command found")
	}

	RequestGitAccess(sshCommand, userId, permissionCheckUrl)
}
