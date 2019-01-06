NOVENDOR := $(shell go list ./... | grep -v vendor)

.PHONY: clean
clean:
	rm -rf ./hello-world/hello-world

.PHONY: build
build:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world

.PHONY: test
test:
	GO111MODULE=on go test -cover -v $(NOVENDOR)
