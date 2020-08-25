IMAGE=skwrl/silicon-dawn:latest
SRV_PKG=./cmd/silicon-dawn
SRV_BIN=./bin/silicon-dawn
DL_PKG=./cmd/download
DL_BIN=./bin/download

all: update test docker-run

update:
	go get -u ./...
	go mod tidy
	go mod vendor

test:
	go test ./...

build: data
	docker build -t $(IMAGE) .

run: docker-run
	 docker run -p 8080:3200 --name Make-Dawn $(IMAGE)

push: build
	docker push $(IMAGE)

local-build: $(BIN)
	go build -v -o $(SRV_BIN) $(SRV_PKG)
	go build -v -o $(DL_BIN) $(DL_PKG)

local: local-build
	$(SRV_BIN)

download: local-build
	$(DL_BIN)
