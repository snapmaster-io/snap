![SnapMaster](https://github.com/snapmaster-io/snapmaster/blob/master/public/SnapMaster-logo-220.png)
# SnapMaster 
## Master your DevOps toolchain

SnapMaster is the definitive DevOps integration platform.  

`snap` is the command-line interface (CLI) for SnapMaster.

## Implementation notes

`snap` is written in golang and communicates with the [SnapMaster-API](https://github.com/snapmaster-io/snapmaster-api) as a back-end.  It utilizes [Cobra](https://github.com/spf13/cobra) for command processing and [Viper](https://github.com/spf13/viper) for config abstraction.

Since SnapMaster currently uses [Auth0](https://auth0.com) for its authentication and authorization, an important part of snap is handling the Proof Key for Code Exchange ([PKCE](https://tools.ietf.org/html/rfc7636)) OAuth2 flow.  

## Source directory structure

### `pkg`
####   `api`: a package that abstracts GET/POST calls against the SnapMaster API
####   `auth`: handle the PKCE authorization flow
####   `cmd`: cobra command implementations
####   `config`: config reading and writing

## Building snap

`go build` in the root directory

## Help

`snap --help` or `snap command --help` makes it easy to learn about all of snap's commands, thanks to Cobra.

## Initializing snap

`snap init` will create a config file (defaults to $HOME/.config/snap/config.json).  This has the most important configuration for snap:

  API URL: the URL for the API.  Currently defaults to https://dev.snapmaster.io
  Client ID: the OAuth2 Client ID for the app.
  Redirect URL: This is the localhost URL where snap expects the OAuth2 callback.  Defaults to http://localhost:8085
  Auth Domain: the OAuth2 server that will handle the PKCE flow. Defaults to snapmaster-dev.auth0.com

`snap init` allows any of these to be overridden.

## Logging in

`snap login` will initiate the login flow.  Note that you must have an account provisioned already on the SnapMaster web app for this to work.  

## Other commands

`snap snaps list` will list all snaps in the user's account

`snap snaps get {snapname}` will get the description of a snap

`snap snaps list --format=json | jq '.[] | .snapId'` will grab the user's snaps in JSON format and pipe through jq, returning a list of the snapId's 