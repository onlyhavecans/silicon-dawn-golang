FROM onlyhavecans.works/oci/golang:1.26@sha256:c05f28d5148bc5c4b60ab5c002291e830b7e835922d23875152b3af5951cecea AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM onlyhavecans.works/oci/static:latest@sha256:c0e338684f4271e71aace102225a72650376f64452cb24135c343a221fa54d3b AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
