FROM onlyhavecans.works/oci/golang:1.25@sha256:93829d2530defccd1184b0bedb7f21fdda98a53d052835ad34667241ba5742d5 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:6ceafbc2a9c566d66448fb1d5381dede2b29200d1916e03f5238a1c437e7d9ea AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
