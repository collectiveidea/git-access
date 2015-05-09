package main

import (
	"fmt"
	shellwords "github.com/mattn/go-shellwords"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

type GitCommandRequest struct {
	FullCommand  string
	CommandParts []string
	BinaryPath   string
	Repository   string
}

// GitRequest takes the requested git command, ensures the user has
// permission to the given repository, and rewrites itself to let
// the git command with it's work (clone or push) to the right repository.
// The permissions request will hit the configured permissionUrl, adding
// two parameters: user and repository. This request needs to return 2xx for
// success, and >= 400 for failure.
func GitRequest(userId string, permissionUrl string) {
	command := parseOriginalCommand()

	if repoAccessAllowed(command, userId, permissionUrl) {
		err := syscall.Exec(command.BinaryPath, command.CommandParts, []string{})
		fmt.Println("Failed to execute the command", command.FullCommand, err)
	}

	os.Exit(1)
}

func parseOriginalCommand() GitCommandRequest {
	fullCommand := os.Getenv("SSH_ORIGINAL_COMMAND")

	if fullCommand == "" {
		fmt.Println("No original command specified. Exiting")
		os.Exit(1)
	}

	commandParts, _ := shellwords.Parse(fullCommand)
	action := commandParts[0]

	if !isValidAction(action) {
		os.Exit(1)
	}

	binaryPath, err := exec.LookPath(action)
	if err != nil {
		fmt.Println("Unable to find the binary", action)
		os.Exit(1)
	}

	var repositoryName string

	if len(commandParts) > 1 {
		repositoryName = commandParts[1]
	} else {
		fmt.Println("Specify the repository.")
		os.Exit(1)
	}

	return GitCommandRequest{
		FullCommand:  fullCommand,
		CommandParts: commandParts,
		BinaryPath:   binaryPath,
		Repository:   repositoryName,
	}
}

func isValidAction(action string) bool {
	return action == "git-receive-pack" ||
		action == "git-upload-pack" ||
		action == "git-upload-archive"
}

func repoAccessAllowed(command GitCommandRequest, userId string, permissionUrl string) bool {
	request, _ := http.NewRequest("GET", permissionUrl, nil)
	values := request.URL.Query()

	values.Add("user", userId)
	values.Add("repository", command.Repository)
	request.URL.RawQuery = values.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Net Error:", err)
		os.Exit(1)
	}

	response.Body.Close()
	return response.StatusCode >= 200 && response.StatusCode < 300
}
