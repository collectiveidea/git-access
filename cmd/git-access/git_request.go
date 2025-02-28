package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"syscall"

	shellwords "github.com/mattn/go-shellwords"
)

type CommandRequest struct {
	command      string
	commandParts []string
	commandPath  string

	user               string
	permissionCheckUrl string
	repository         string
}

func (self *CommandRequest) RewriteRepository(newRepo string) {
	if newRepo != "" {
		self.commandParts[1] = newRepo
		self.repository = newRepo
	}
}

// RequestGitAccess takes the requested git command, ensures the user has
// permission to the given repository, and rewrites itself to let
// the git command with it's work (clone or push) to the right repository.
// The permissions request will hit the configured permissionUrl, adding
// two parameters: user and repository. This request needs to return 2xx for
// allowed, and >= 400 for denied.
//
// This permissions request can also return in the body of the response the local
// path to the real repository on disk, in which the command will be rewritten
// to point to the actual repository before being exec'd.
func RequestGitAccess(gitCommand string, userId string, permissionCheckUrl string) {
	request := validateRequest(gitCommand, userId, permissionCheckUrl)

	if repoAccessAllowed(&request) {
		log.Fatal(
			"Failed to execute command.",
			executeOriginalRequest(&request),
		)
	} else {
		log.Fatalf("Permission denied.")
	}
}

func validateRequest(command string, userId string, permissionCheckUrl string) (request CommandRequest) {
	commandParts, _ := shellwords.Parse(command)
	binary := commandParts[0]

	if !isValidAction(binary) {
		log.Fatalf("Permission denied.")
	}

	binaryPath, err := exec.LookPath(binary)
	if err != nil {
		log.Fatal("Unknown command.", binary)
	}

	var repository string
	if len(commandParts) > 1 {
		repository = commandParts[1]
	} else {
		log.Fatalf("Missing repository.")
	}

	request = CommandRequest{
		command:            command,
		commandParts:       commandParts,
		commandPath:        binaryPath,
		user:               userId,
		permissionCheckUrl: permissionCheckUrl,
		repository:         repository,
	}

	return
}

func isValidAction(action string) bool {
	return action == "git-receive-pack" ||
		action == "git-upload-pack" ||
		action == "git-upload-archive"
}

func repoAccessAllowed(request *CommandRequest) bool {
	permissionCheck, _ := http.NewRequest("GET", request.permissionCheckUrl, nil)
	values := permissionCheck.URL.Query()

	values.Add("user", request.user)
	values.Add("repository", request.repository)
	permissionCheck.URL.RawQuery = values.Encode()

	response, err := http.DefaultClient.Do(permissionCheck)
	if err != nil {
		log.Fatalf("Net Error: %v\n", err)
	}
	defer response.Body.Close()

	responseSuccess := response.StatusCode >= 200 && response.StatusCode < 300

	body, _ := io.ReadAll(response.Body)
	request.RewriteRepository(strings.TrimSpace(string(body)))

	return responseSuccess
}

func executeOriginalRequest(request *CommandRequest) error {
	return syscall.Exec(request.commandPath, request.commandParts, []string{})
}
