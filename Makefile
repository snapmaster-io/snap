# grab the build hash from the latest git commit hash
BUILD_HASH=$(shell git rev-parse HEAD)

# inject the build hash into the GitHash global in the version package
# also use -s -w to strip symbols out of the build
LDFLAGS="-s -w -X github.com/snapmaster-io/snap/pkg/version.GitHash=${BUILD_HASH}"

all:
	go build -o bin/snap -ldflags $(LDFLAGS)
