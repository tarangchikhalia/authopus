# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/authopus ./cmd/authopus

FROM scratch
COPY --from=build /out/authopus /authopus
USER 65532:65532
ENTRYPOINT ["/authopus"]
