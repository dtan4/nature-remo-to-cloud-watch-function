.PHONY: generate
generate: mockgen
	go generate -v ./...

.PHONY: mockgen
mockgen:
	go install -v github.com/golang/mock/mockgen

.PHONY: test
test:
	GO111MODULE=on go test -coverpkg=./... -coverprofile=coverage.txt -v ./...
