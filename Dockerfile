FROM onlyhavecans.works/oci/golang:1.25@sha256:b3dde9df96709c60c299573912ee80d5e91dc28a8021015e620276a017ec86e7 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:b55f6779fb7990fb7db5e272c69a4cd6ea7070f3195da71b5ae163bfdbef4f76 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
