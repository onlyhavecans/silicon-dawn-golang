FROM onlyhavecans.works/oci/golang:1.25@sha256:f7c6f0bf2d22b0b840d2b4956466b155532de36351056eaba6ea5e1aa458dc3b AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:838610585824d0141daf5d76af5778f59a6f9dcc1a822148c790c9043c89e8eb AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
