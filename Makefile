NOVENDOR := $(shell go list ./... | grep -v vendor)

.PHONY: test
test:
	GO111MODULE=on go test -cover -v $(NOVENDOR)
