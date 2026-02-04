FROM onlyhavecans.works/oci/golang:1.25@sha256:96323c4aa0ea9064c4a4ac0cee942c235173d2674daa641cccbbc021fec18b6a AS build

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
