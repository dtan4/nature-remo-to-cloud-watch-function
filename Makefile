.PHONY: generate
generate: mockgen
	go generate -v ./...

.PHONY: mockgen
mockgen:
	go install -v github.com/golang/mock/mockgen

.PHONY: test
test:
	GO111MODULE=on go test -coverprofile=coverage.txt -v `go list ./... | grep -v aws/mock`
