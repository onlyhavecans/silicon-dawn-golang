FROM onlyhavecans.works/oci/golang:1.25@sha256:4c22673170524e09a47a6e1a5f1ca99692f4796c9b4f2de25ec7137bcd897edd AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:6f3f2123de90d2e7998b8161a2838433ec32560a827d07bcab339dacbf0cf16f AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
