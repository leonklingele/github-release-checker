# Github Release checker

[![Build Status](https://travis-ci.org/leonklingele/github-release-checker.svg?branch=master)](https://travis-ci.org/leonklingele/github-release-checker)

Get notified on new releases of your starred & watched repos.

## Setup

```sh
# go get this repo
go get -u github.com/leonklingele/github-release-checker
# .. and cd into it
cd $GOPATH/src/github.com/leonklingele/github-release-checker

# Install config file to $HOME/.github-release-checker/config.toml
make config
# .. and edit it at will
# I recommend to not enable "mail" on the first run as it will
# most likely spam your inbox.
$EDITOR $HOME/.github-release-checker/config.toml

# Finally start the app
./github-release-checker
```

## Example config

My current configuration:

```toml
[checker]
# How frequently to run the check
interval = "5m"

[checker.db]
path = "$HOME/.github-release-checker/sqlite.db"

[checker.repositories]
# Repos to ignore
ignored   = [
  "vim/vim",
]
# Repos to mark as "important". This is useful to e.g. send push
# notification emails for releases of certain repos.
important = [
  "git/git",
  "golang/go",
  "libressl-portable/openbsd",
  "libressl-portable/portable",
  "nginx/nginx",
  "openssh/openssh-portable",
  "openssl/openssl",
  "openvpn/openvpn",
]

[checker.github]
# Your Github username
user  = "leonklingele"
# Your Github access token. Only "public access" must be granted.
# Generate a token here: https://github.com/settings/tokens
token = ".."

[mail]
# Whether to send email notifications on new releases
enabled  = true
# Number of mail workers
workers  = 5
# Whether to accept untrusted certificates
insecure = true

# How to connect to the mail server
host = "localhost"
port = 25
user = ""
pswd = ""

# Sender address
from         = "github-releases@leonklingele.de"
# Recipient addresses
to           = [ "github-releases@leonklingele.de" ]
# Recipient addresses for "important" releases, e.g. a Boxcar
# email address: https://boxcar.io/
important_to = [ "boxcar-push@leonklingele.de" ]

# Subject template
subject = "New release of $fullName"
# Body template
body    = "New release of $url : $version"
```
