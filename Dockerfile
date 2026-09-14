FROM onlyhavecans.works/oci/golang:1.27@sha256:b475798fb16158e6c38e8b5ca2d870fbeaa8b7fec0fc8ec64b3dc20966040635 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:a8a6666cb7b13df7d18fd65b4264c2edafbce92acc0f32c423b5470d7d717665 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
