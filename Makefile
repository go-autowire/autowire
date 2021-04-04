# Make

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint: $(GOLINT)
	golangci-lint run --fix --fast ${FileDir} -e "(\\w|\\s|\`)+\\w+(\`)? should be(\\w|\\s|\`)+" --skip-dirs=fake

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html