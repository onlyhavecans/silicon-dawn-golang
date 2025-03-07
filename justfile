# Variables
image := "skwrl/silicon-dawn:latest"
srv_pkg := "./cmd/silicon-dawn"
srv_bin := "./bin/silicon-dawn"
dawnzip := "The-Tarot-of-the-Silicon-Dawn.zip"
cards := "data"

# Default recipe
default:
  @just --list

# Run lint, fmt, test, and docker-run
all: lint fmt test docker-run

# Update dependencies
update:
    go get -u ./...
    go mod tidy
    go mod vendor
    git diff

# Run linter
lint:
    golangci-lint run

# Format code
fmt:
    go install mvdan.cc/gofumpt@latest
    go fmt ./...
    gofumpt -w ./

# Run tests
test:
    go vet ./...
    go test ./...

# Build Docker image
build: ensure-cards
    docker build -t {{image}} .

# Run in Docker container
docker-run: build
    docker run --rm -p 8080:3200 --name Make-Dawn {{image}}

# Push Docker image
push: build
    docker push {{image}}

# Build locally
local-build:
    go build -v -o {{srv_bin}} {{srv_pkg}}

# Run locally
local: local-build
    {{srv_bin}}

# Download tarot zip file
download-zip:
    wget "http://egypt.urnash.com/media/blogs.dir/1/files/2018/01/The-Tarot-of-the-Silicon-Dawn.zip"

# Extract cards
extract-cards: download-zip
    mkdir -p {{cards}}
    unzip -oj {{dawnzip}} -x "__MACOSX/*" "*/sand-home*" -d {{cards}}

# Ensure cards directory exists
ensure-cards:
    #!/usr/bin/env bash
    if [ ! -d "{{cards}}" ] || [ -z "$(ls -A {{cards}})" ]; then
        just download-zip
        just extract-cards
    fi

# Clean built artifacts
clean:
    rm -f {{srv_bin}}
    docker rmi {{image}} || true

# Run a specific test by name
test-one test_name:
    go test -v ./... -run "{{test_name}}"
