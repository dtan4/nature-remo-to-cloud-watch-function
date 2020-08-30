NAME := nature-remo-to-cloud-watch-function

LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := build

.PHONY: build
build:
	cd function; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o ../bin/$(NAME)

.PHONY: test
test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...
