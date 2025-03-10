package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
//	[
//	  { user: "1", keys: ["ssh-rsa AAA...==", "ssh-rsa AAB...=="]},
//	  { user: "2", keys: ["ssh-rsa AAD...=="]},
//	  ...
//	]
func RequestAuthorizedKeys(commandBinary string, keysUrl string) {
	users := readKeys(keysUrl)

	for _, user := range users {
		for _, publicKey := range user.Keys {
			fmt.Printf(
				"command=\"%s --user=%d\",%s %s\n",
				commandBinary,
				user.UserId,
				AuthorizedKeysOptions,
				publicKey,
			)
		}
	}
}

func readKeys(url string) (keysList []UserKeys) {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal("Error receiving keys", err)
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&keysList)

	if err != nil {
		log.Fatal("Error parsing keys response", err)
	}

	return
}
