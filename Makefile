IMAGE=skwrl/silicon-dawn:latest
SRV_PKG=./cmd/silicon-dawn
SRV_BIN=./bin/silicon-dawn
DAWNZIP=The-Tarot-of-the-Silicon-Dawn.zip
CARDS=data

all: lint fmt test docker-run

update:
	go get -u ./...
	go mod tidy
	go mod vendor
	git diff

lint:
	golangci-lint run

fmt:
	go install mvdan.cc/gofumpt@latest
	go fmt ./...
	gofumpt -w ./

test:
	go test ./...

build: $(CARDS)
	docker build -t $(IMAGE) .

docker-run: build
	docker run --rm -p 8080:3200 --name Make-Dawn $(IMAGE)

push: build
	docker push $(IMAGE)

local-build: $(BIN)
	go build -v -o $(SRV_BIN) $(SRV_PKG)

local: local-build
	$(SRV_BIN)

$(DAWNZIP):
	wget "http://egypt.urnash.com/media/blogs.dir/1/files/2018/01/The-Tarot-of-the-Silicon-Dawn.zip"

$(CARDS): $(DAWNZIP)
	unzip -oj $(DAWNZIP) -x "__MACOSX/*" "*/sand-home*" -d $(CARDS)
