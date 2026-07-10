FROM onlyhavecans.works/oci/golang:1.25@sha256:09c6d487ccb96cac78767ef217cef33d15e9ee8c7569edbc9a3b00e3aef505d5 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:89a7f06296db723064812805b50f16717e8e4150cdd883e89378e05f410a7b9d AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
