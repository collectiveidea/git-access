package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
func RequestAuthorizedKeys(commandBinary string, keysUrl string) error {
	users, err := readKeys(keysUrl)
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, publicKey := range user.Keys {
			fmt.Println(
				"command=\""+commandBinary+" --user="+strconv.Itoa(user.UserId)+"\","+AuthorizedKeysOptions,
				publicKey,
			)
		}
	}

	return nil
}

func readKeys(url string) (keysList []UserKeys, err error) {
	response, err := http.Get(url)

	if err != nil {
		err = fmt.Errorf("Error receiving keys", err)
		return
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &keysList)

	if err != nil {
		err = fmt.Errorf("Error parsing response", err)
	}

	return
}
