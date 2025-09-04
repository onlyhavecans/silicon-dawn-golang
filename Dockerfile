FROM golang:1.25@sha256:5e856b892fff06368bd92fe3ac00c457ede6d399e7152575b57f9a14312ab713 AS build

ENV GOFLAGS="-mod=vendor"

WORKDIR /go/src/app
COPY . .

RUN go vet ./... && go test ./...

RUN CGO_ENABLED=0 go install -trimpath ./cmd/silicon-dawn

# Final Stage
# FROM scratch AS production
FROM gcr.io/distroless/static:latest@sha256:f2ff10a709b0fd153997059b698ada702e4870745b6077eff03a5f4850ca91b6 AS production
EXPOSE 3200/tcp

COPY --from=build /go/bin/silicon-dawn /
COPY templates /templates
COPY data /data

USER nonroot
CMD ["/silicon-dawn"]
