# Variables
image := "skwrl/silicon-dawn:latest"
srv_pkg := "./cmd/silicon-dawn"
srv_bin := "./bin/silicon-dawn"
dawnzip := "The-Tarot-of-the-Silicon-Dawn.zip"
cards := "data"

# Default recipe
default:
    @just --list --unsorted

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
    golangci-lint fmt

# Autofix lint issues. --disable govet: its fieldalignment autofix reorders
# struct fields without updating positional literals elsewhere — silent
# corruption, not a lint nit. fieldalignment is still reported by `just lint`.
fix:
    go fix ./...
    golangci-lint run --fix --disable govet

# Run tests
test:
    go vet ./...
    go test ./...

# Run lint and test
check: lint test

# Build Docker image
build: ensure-cards
    docker build -t {{ image }} .

# Run in Docker container
docker-run: build
    docker run --rm -p 8080:3200 --name Make-Dawn {{ image }}

# Push Docker image (legacy path; releases now go through GoReleaser, see `release`)
push: build
    docker push {{ image }}

# Build local snapshot images via GoReleaser (no push, for testing dockers_v2 config)
release-snapshot: ensure-cards
    goreleaser release --snapshot --clean

# Validate .goreleaser.yaml
goreleaser-check:
    goreleaser check

# Cut a new release: tags vX.Y.Z and pushes it, triggering the release workflow
release version:
    #!/usr/bin/env bash
    set -e
    if [ -n "$(git status --porcelain)" ]; then
        echo "working tree is dirty — commit or stash first" >&2
        exit 1
    fi
    git tag -a "v{{ version }}" -m "release v{{ version }}"
    git push origin "v{{ version }}"

# Build locally
local-build:
    go build -v -o {{ srv_bin }} {{ srv_pkg }}

# Run locally
local: local-build
    {{ srv_bin }}

# Download tarot zip file
download-zip:
    wget "http://egypt.urnash.com/media/blogs.dir/1/files/2018/01/The-Tarot-of-the-Silicon-Dawn.zip"

# Extract cards
extract-cards: download-zip
    mkdir -p {{ cards }}
    unzip -oj {{ dawnzip }} -x "__MACOSX/*" "*/sand-home*" -d {{ cards }}

# Ensure cards directory exists
ensure-cards:
    #!/usr/bin/env bash
    if [ ! -d "{{ cards }}" ] || [ -z "$(ls -A {{ cards }})" ]; then
        just download-zip
        just extract-cards
    fi

# Clean built artifacts
clean:
    rm -f {{ srv_bin }}
    docker rmi {{ image }} || true

# Run a specific test by name
test-one test_name:
    go test -v ./... -run "{{ test_name }}"
