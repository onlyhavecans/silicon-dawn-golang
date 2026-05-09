FROM onlyhavecans.works/oci/golang:1.25@sha256:c3b7d08caebaf7de38d2640411e48bdc828d51add8f8f182e4ed163cff955f98 AS build

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
