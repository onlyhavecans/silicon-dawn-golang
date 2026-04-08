FROM onlyhavecans.works/oci/golang:1.25@sha256:ed520ab5bed37ce887012c050ced60f7d52dfcd212e3dc6fdd8951e9c4e25c1a AS build

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
