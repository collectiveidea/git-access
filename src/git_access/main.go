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
			keysUrl := c.String("authorized-keys-url")

			if keysUrl == "" {
				fmt.Println("Error: --authorized-keys-url is required when --authorized-keys is used")
				os.Exit(1)
			}

			AuthorizedKeys(keysUrl)
		} else {
			GitRequest(c)
		}
	}

	app.Run(os.Args)
}
