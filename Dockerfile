FROM onlyhavecans.works/oci/golang:1.25@sha256:7af63db8d8dc56289c8fa6d9883ad9d043c332755343a243dbb5d91984343a03 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:227aa7b4f3d89833db58676eacdbe9a49b5d5e4748e0ec3f05005335fa73aaf9 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
