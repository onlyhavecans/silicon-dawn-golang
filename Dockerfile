FROM golang:1.15 AS builder

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

ARG CGO_ENABLED=0
RUN go install ./cmd/silicon-dawn

# Final Stage
FROM scratch AS production
EXPOSE 3200/tcp

COPY --from=builder /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

CMD ["./silicon-dawn"]
