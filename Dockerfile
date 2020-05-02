FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /data
ENV CardsDirectory /data
EXPOSE 3200/tcp

CMD ["silicondawn", "serve", "--release"]