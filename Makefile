BUILD_HASH = $(shell git rev-parse HEAD)

all:
	go build -o bin/snap -ldflags "-X github.com/snapmaster-io/snap/pkg/version.GitHash=${BUILD_HASH}"
