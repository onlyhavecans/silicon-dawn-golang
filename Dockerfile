FROM onlyhavecans.works/oci/golang:1.25@sha256:fb4095b65a7bb89f039def7e33d7b90095d2c25f34597748758a6f209eead7ff AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:eb5a533892990d7a3ad187778d6274054bda0cac40491129824a0c035b6e8fb3 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
