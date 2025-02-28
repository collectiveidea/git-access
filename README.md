# git-access

Provide remote access to git repositories through SSH and protection via key pairs.

`git-access` solves a very specific set of problems. If you need:

* To host git repositories for others
* Make `git clone` available over SSH
* Restrict who can clone repositories with SSH public keys

then `git-access` is for you!

**Note:** `git-access` requires OpenSSL 6.2 or later as it takes advantage of the `AuthorizedKeysCommand` and `AuthorizedKeysCommandUser` options added in that version.

Running `git clone` over SSH is a two step process. First, SSH needs to confirm that the incoming public key is allowed to access the box, and second the application in question needs to confirm that the user in question has access to the requested repository. `git-access` encompases both steps in one binary.

## Usage

`git-access` has two modes of operation: Authorization and Access.

### Authorization

The authorization step is when the user's SSH keys are checked against the known list of keys from the application. `git-access` must be configured with a URL from which the list of valid SSH keys will be returned. The response must be a JSON array of objects containing two keys: `user` and `keys`. This request must also return *all* users and their SSH keys as the selection and authorization logic is handled inside of SSH itself.

* `user` :: String :: Unique identifier the application gives user accounts
* `keys` :: Array  :: An array of public keys for that user

```json
[
  {"user": "1", "keys": ["ssh-rsa AAA...==", "ssh-rsa AAB..=="]},
  {"user": "2", "keys": ["ssh-rsa AAD...=="]},
]
```

The options supported in this mode are:

```
git-access \
  --authorized-keys         # Enable Authorization mode. Also available as "-A"
  --authorized-keys-url=    # Full URL that returns SSH keys in the format above
  --authorize-command=      # Path to binary for initiating the actual git access
                            # This will default to the current `git-access` binary
                            # but for reasons stated below may be better to be explicit here.
```

### Access

Once SSH authorization completes, it is time to ensure the current user has access to the repository in question. This consists of a second HTTP request which is expected to return with a status code of 2xx for success and >= 400 for failure or denied. This request will receive two parameters: `user` and `repository`, where `user` will match the identifier returned from Authorization and `repository` will be pulled from the original `git clone` request. E.g. for `git clone git@my-app.com:my-site/blog.git`, `repository` will equal `"my-site/blog.git"`.

The options supported in this mode are:

```
git-access \
  --user=                  # The unique user identifier, usually provided from Authorization
  --permission-check-url=  # The URL which will check if the user has access or not.
```

### Global Options

git-access also supports the following flags regardless of execution mode:

```
git-access \
  --syslog                 # Enable logging to a local syslog daemon. Will log under info as "git-access".
```

## Server Configuration

The following is the recommended way of configuring `git-access` on servers. The `AuthorizedKeysCommand` works best when it is given a single binary file to run, thus it is recommended to set up shell scripts that call out to `git-access` with the required parameters.  You'll need two of these such scripts:

#### get-authorized-keys.sh

This script will call out to get the list of known keys and will look something like:

```sh
#!/bin/sh

/path/to/git-access -A --authorized-keys-url=http://your.app.com/git_access/keys --authorize-command=/path/to/confirm-git-access.sh
```

#### confirm-git-access.sh

Once SSH keys are confirmed, checking repository access permissions and forwarding along the original git request are the second script. The `get-authorized-keys.sh` command will inject `--user` into this script depending on what came back from the server during authorization.

```sh
#!/bin/sh

/path/to/git-access "$@" --permission-check-url=http://your.app.com/git_access/access
```

#### sshd_config

Then, to tie all of the above together, configure SSH itself to run the first script.

**Note:** The script called by `AuthorizedKeysCommand` **must** be owned by root and **must** be writable only by root, or SSH will not execute the script. Also this requires `PubkeyAuthentication` to be turned on.

```
Match User git
  AuthorizedKeysCommand /path/to/get-authorized-keys.sh
  AuthorizedKeysCommandUser git
```

Restart SSH to apply the new configuration and enjoy!

## Developing

After cloning this repository, you can build `git-access` with `go build ./cmd/git-access` or `rake build`.

Due to the nature of the system, the tests are high level that set up actual web servers and check that the full communication process works. As such, the easiest way I know to do this is through Ruby, so the tests require Ruby 2.0 or greater. Run tests with `rake`.

## Contributing

`git-access` is open source and contributions are encouraged! No contribution is too small.

Please see the [contribution guidelines](CONTRIBUTING.md) for more information.

