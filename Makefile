NAME := nature-remo-to-cloud-watch-function

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	docker-compose run --rm -e CGO_ENABLED=0 go build $(LDFLAGS) -a -tags netgo -installsuffix netgo -o bin/$(NAME) github.com/dtan4/nature-remo-to-cloud-watch-function/function

.PHONY: deploy
deploy: bin/$(NAME)
ifeq ($(AWS_S3_BUCKET),)
	@echo "AWS_S3_BUCKET must be set" >&2
	@exit 1
endif
ifeq ($(AWS_CLOUDFORMATION_STACK_NAME),)
	@echo "AWS_CLOUDFORMATION_STACK_NAME must be set" >&2
	@exit 1
endif
	docker-compose run --rm sam package --template-file template.yaml --s3-bucket $(AWS_S3_BUCKET) --output-template-file packaged.yaml
	docker-compose run --rm sam deploy --template-file packaged.yaml --stack-name $(AWS_CLOUDFORMATION_STACK_NAME) --capabilities CAPABILITY_IAM

.PHONY: setup
setup: setup-go setup-sam

.PHONY: setup-go
setup-go:
	docker-compose build go

.PHONY: setup-sam
setup-sam:
	docker-compose build sam

.PHONY: test
test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...
