#--- Build stage
FROM golang:1.24-bullseye AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0

#--- Image stage
FROM alpine:3.22.0

COPY --from=go-builder /src/target/axone-mcp /usr/bin/axone-mcp

ENTRYPOINT ["/usr/bin/axone-mcp"]
