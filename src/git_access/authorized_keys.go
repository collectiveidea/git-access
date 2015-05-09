package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	AuthorizedKeysOptions = "no-user-rc,no-X11-forwarding,no-agent-forwarding,no-pty"
)

type UserKeys struct {
	UserId int      `json:"user_id"`
	Keys   []string `json:"keys"`
}

// AuthorizedKeys queries the given keysUrl for a list of known SSH public keys.
// These keys are transformed into the authorized_keys file format, configured
// such that this tool can turn around and ensure the requesting user has permission
// to the repository they are requesting.
//
// The response is expected to be a JSON array with each entry including the user_id
// and a list of keys for that user:
//
//   [
//     { user_id: 1, keys: ["ssh-rsa AAA...==", "ssh-rsa AAB...=="]},
//     { user_id: 2, keys: ["ssh-rsa AAD...=="]},
//     ...
//   ]
//
func AuthorizedKeys(keysUrl string) {
	for _, user := range readKeys(keysUrl) {
		for _, publicKey := range user.Keys {
			fmt.Println(
				"command=\"git-access --user="+strconv.Itoa(user.UserId)+"\","+AuthorizedKeysOptions,
				publicKey,
			)
		}
	}
}

func readKeys(url string) (keysList []UserKeys) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Net Error:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	if err = json.Unmarshal(responseBody, &keysList); err != nil {
		fmt.Println("Error parsing JSON response", err)
		os.Exit(1)
	}

	return
}
