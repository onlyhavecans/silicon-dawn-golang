FROM onlyhavecans.works/oci/golang:1.27@sha256:31df75a0f51705bb15c74cacaeadb6596ae270cea9ec05138928de7e2d1f65e8 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:0985f124d25d79a432b79e806764a9deb759e5c664be7c0633b9f13c3e12cbc0 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
