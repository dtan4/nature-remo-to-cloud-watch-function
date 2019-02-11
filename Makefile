NAME := nature-remo-to-cloud-watch-function

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	docker-compose run --rm -e CGO_ENABLED=0 go build $(LDFLAGS) -a -tags netgo -installsuffix netgo -o bin/$(NAME) github.com/dtan4/nature-remo-to-cloud-watch-function/function

.PHONY: generate
generate:
	docker-compose run --rm go generate -v ./...

.PHONY: test
test:
	docker-compose run --rm go test -coverprofile=coverage.txt -v `docker-compose run -T --rm go list ./... | grep -v aws/mock`
