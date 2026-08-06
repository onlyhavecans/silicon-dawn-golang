# Variables
image := "skwrl/silicon-dawn:latest"
srv_pkg := "./cmd/silicon-dawn"
srv_bin := "./bin/silicon-dawn"
dawnzip := "The-Tarot-of-the-Silicon-Dawn.zip"
cards := "data"

# List available recipes.
default:
    @just --list

# Build local binaries into ./bin/.
build:
    mkdir -p bin
    go build -trimpath -o bin/ ./cmd/silicon-dawn

# Run the test suite with the race detector.
test:
    go test -race ./...

# Lint.
lint:
    golangci-lint run

# Format all Go source.
fmt:
    golangci-lint fmt

# Lint with autofix.
fix:
    golangci-lint run --fix

# Lint + test.
check: lint test

# Tidy go.mod/go.sum.
tidy:
    go mod tidy
    go mod vendor

# Update all dependencies and tidy.
update:
    go get -u ./...
    go mod tidy
    go mod vendor

# Remove build artifacts.
clean:
    rm -rf bin/ dist/

# Repo-specific recipes

# Build Docker image
docker-build: ensure-cards
    docker build -t {{ image }} .

# Run in Docker container
docker-run: docker-build
    docker run --rm -p 8080:3200 --name Make-Dawn {{ image }}

# Push Docker image
docker-push: docker-build
    docker push {{ image }}

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

# Cut and push a release tag (accepts 1.2.3 or v1.2.3).
release version: check
    #!/usr/bin/env bash
    set -euo pipefail
    ver="{{ version }}"
    ver="v${ver#v}"
    [[ "$ver" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "invalid version '$ver' (want X.Y.Z)" >&2; exit 1; }
    [[ -z "$(git status --porcelain)" ]] || { echo "working tree not clean" >&2; exit 1; }
    [[ "$(git branch --show-current)" == "main" ]] || { echo "releases cut from main only" >&2; exit 1; }
    git fetch --quiet --tags origin main
    [[ "$(git rev-parse HEAD)" == "$(git rev-parse origin/main)" ]] || { echo "main not in sync with origin" >&2; exit 1; }
    git tag -a "$ver" -m "$ver"
    git push origin "$ver"
    echo "tagged $ver — release workflow is running (fj actions tasks)"

# Tag the next major/minor/patch release.
release-major: (release `just _next major`)
release-minor: (release `just _next minor`)
release-patch: (release `just _next patch`)

# Compute the next version from the latest tag (v-prefixed or bare).
_next kind:
    #!/usr/bin/env bash
    set -euo pipefail
    git fetch --quiet --tags origin 2>/dev/null || true
    latest=$(git tag --list 'v[0-9]*' '[0-9]*' | sed 's/^v//' | sort -V | tail -1)
    IFS=. read -r maj min pat <<<"${latest:-0.0.0}"
    case "{{ kind }}" in
      major) echo "v$((maj+1)).0.0" ;;
      minor) echo "v${maj}.$((min+1)).0" ;;
      patch) echo "v${maj}.${min}.$((pat+1))" ;;
    esac
