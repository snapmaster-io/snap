![SnapMaster](https://github.com/snapmaster-io/snapmaster/blob/master/public/SnapMaster-logo-220.png)
# SnapMaster 
## Master your DevOps toolchain

SnapMaster is the definitive DevOps integration platform.  

`snap` is the command-line interface (CLI) for SnapMaster.

## Installing or building `snap`

### Installing the latest version

On a Mac: `brew tap snapmaster-io/snap && brew install snap`

If you have go installed: `go get github.com/snapmaster-io/snap`

### Building `snap` from source

`make` in the root directory will invoke `go build -o bin/snap` and embed the latest git hash into the version.

## Using `snap`

### Help

`snap --help` or `snap command --help` makes it easy to learn about all of snap's commands, thanks to Cobra.

### Initializing snap

`snap init` will create a config file (defaults to $HOME/.config/snap/config.json).  This has the most important configuration for snap:

* API URL: the URL for the API.  Currently defaults to https://www.snapmaster.io
* Client ID: the OAuth2 Client ID for the app.
* Redirect URL: This is the localhost URL where snap expects the OAuth2 callback.  Defaults to http://localhost:8085
* Auth Domain: the OAuth2 server that will handle the PKCE flow. Defaults to snapmaster-dev.auth0.com

`snap init` allows any of these to be overridden.

### Logging in

`snap login` will initiate the login flow.  If you don't have a SnapMaster 
account, you can create one.  

`snap logout` will remove the API access token and log out the current user.

### Snap management

#### Interacting with the Gallery

`snap gallery` will retrieve all the snaps in the gallery

`snap gallery get {snapname}` will get the YAML description of a snap

`snap snaps fork {snapname}` will fork a public snap into the user's account

#### Managing your own snaps

`snap snaps list` will list all snaps in the user's account

`snap snaps get {snapname}` will get the YAML description of a snap

`snap snaps list --format=json | jq '.[] | .snapId'` will grab the user's snaps in JSON format and pipe through jq, returning a list of the snapId's 

`snap snaps delete {snapname}` will delete a snap from the user's account

`snap snaps publish/unpublish {snapname}` will make a snap public (discoverable) or switch it back to private

#### Activating and managing active snaps

`snap activate {snapname}` will prompt for parameters and activate a snap

`snap active list` will list all activated snaps 

`snap active get {active snap ID}` will get information about the active snap

`snap active logs {active snap ID}` will get all logs for the active snap

`snap active logs {active snap ID} details {log ID}` will retrieve log details for a particular log entry

`snap active pause/resume {active snap ID}` will pause or resume an active snap

`snap active deactivate {active snap ID}` will deactivate and active snap and REMOVE ALL LOGS

#### Interacting with logs

`snap logs` will retrieve all logs from all active snaps

`snap logs details {logID}` will retrieve log details for a particular log entry

## Source directory structure

### `pkg`
####   `api`: a package that abstracts GET/POST calls against the SnapMaster API
####   `auth`: handle the PKCE authorization flow
####   `cmd`: cobra command implementations
####   `config`: config reading and writing
####   `print`: printing out API responses in all supported formats for all API's
####   `utils`: color-printing support and other generic utilities
####   `version`: version information, with an injectable git hash

## Implementation notes

`snap` is written in golang and communicates with the [SnapMaster-API](https://github.com/snapmaster-io/snapmaster-api) as a back-end.  It utilizes [Cobra](https://github.com/spf13/cobra) for command processing and [Viper](https://github.com/spf13/viper) for config abstraction.

Since SnapMaster currently uses [Auth0](https://auth0.com) for its authentication and authorization, an important part of snap is handling the Proof Key for Code Exchange ([PKCE](https://tools.ietf.org/html/rfc7636)) OAuth2 flow.  The `auth` package provides an implementation of this flow.

