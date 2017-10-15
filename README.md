# Github Release checker

Get notified on new releases of your starred & watched repos.

## Setup

```sh
# Clone this repo
git clone https://github.com/leonklingele/github-release-checker
# .. and cd into it
cd github-release-checker

# Build the app
make
# Create a config file
cp config.toml.in config.toml
# .. and edit it at will
# I recommend to not enable "mail" on the first run as it will
# most likely spam your inbox.
$EDITOR config.toml

# Finally start the app
./github-release-checker
```

## Example config

My current configuration:

```toml
[checker]
# How frequently to run the check
interval = "5m"

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
# Whether to send emails notifications on new releases
enabled  = false
# Number of mail workers
workers  = 5
# Whether to accept untrusted certificates
insecure = false

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
