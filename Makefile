.PHONY: all build test clean coverage docker docker-compose-up docker-compose-down

# Build variables
BINARY_NAME=event-tracking
MAIN_PACKAGE=./cmd
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

all: test build

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

test:
	go test -v -race ./...

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Docker targets
docker:
	docker compose build

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down

# Development tools installation
.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	sudo pacman -S pre-commit 
	

# Initialize development environment
.PHONY: init
init: install-tools
	pre-commit install
	go mod tidy

# Run the application
.PHONY: run
run:
	go run $(MAIN_PACKAGE)

# Generate mocks (requires mockgen)
.PHONY: generate-mocks
generate-mocks:
	go generate ./... 