package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	AuthorizedKeysOptions = "no-user-rc,no-X11-forwarding,no-agent-forwarding,no-pty"
)

// AuthorizedKeys queries the given keysUrl for a list of known SSH public keys.
// These keys are transformed into the authorized_keys file format, configured
// such that this tool can turn around and ensure the requesting user has permission
// to the repository they are requesting.
//
// The response is expected to be a new-line seperated text file with a user identifier
// prepended to the key, seperated by a comma. E.g.
//
//   1,ssh-rsa AAA...==
//   1,ssh-dsa AAB...==
//   2,ssh-rsa AAC...==
//   ...
//
func AuthorizedKeys(keysUrl string) {
	allKeys := readKeys(keysUrl)

	var parts []string
	var userId, userKey string
	for _, key := range allKeys {
		parts = strings.SplitN(key, ",", 2)
		userId = parts[0]
		userKey = parts[1]

		fmt.Println(
			"command=\"git-access --user="+userId+"\","+AuthorizedKeysOptions,
			userKey,
		)
	}
}

func readKeys(url string) []string {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Net Error:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	keys, _ := ioutil.ReadAll(response.Body)
	return strings.Split(string(keys), "\n")
}
