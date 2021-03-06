# Make

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint: $(GOLINT)
	golangci-lint run --fix --fast ${FileDir}

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html