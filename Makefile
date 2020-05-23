IMAGE=skwrl/silicon-dawn:latest
BIN=./bin/silicon-dawn

all: update test docker-run

update:
	go get -u ./...
	go mod tidy
	go mod vendor

test:
	go test ./...

build:
	docker build -f localcards.dockerfile -t $(IMAGE) .

run: docker-build
	 docker run -p 8080:3200 --name Make-Dawn $(IMAGE)

push:
	docker push $(IMAGE)

local-build: $(BIN)
	go build -v -o $(BIN)

local: local-build
	$(BIN) serve

download: local-build
	$(BIN) get

