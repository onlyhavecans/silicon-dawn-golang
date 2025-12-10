FROM onlyhavecans.works/oci/golang:1.25@sha256:828328b1b4b24d4ed279e22e8585f6f4f54af62404d484a0e71eadeb2c18efd3 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:14d584085808e8b2d8c6f24537694a35cc87f7cbee39493f3fa3cee0d2eef13c AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
