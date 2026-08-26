FROM onlyhavecans.works/oci/golang:1.27@sha256:eef6a67266eeed3c86dd47fd01b32faa8bf0229eb83eb3d4d466e80391bd3820 AS build

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
