FROM golang:1.14 AS builder

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

ARG CGO_ENABLED=0
RUN go install ./...

# Final Stage
FROM scratch
EXPOSE 3200/tcp

COPY --from=builder /go/bin/silicondawn /
COPY templates /templates

COPY data /data

ENTRYPOINT ["./silicondawn"]
CMD ["serve", "--release"]
