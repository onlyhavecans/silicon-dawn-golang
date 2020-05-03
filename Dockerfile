FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
ARG CGO_ENABLED=0
RUN go install -v ./...


FROM scratch

COPY --from=builder /go/bin/silicondawn /
COPY templates /templates

# RUN ["silicondawn", "get"]
COPY data /data

EXPOSE 3200/tcp
ENTRYPOINT ["./silicondawn"]
CMD ["serve", "--release"]
