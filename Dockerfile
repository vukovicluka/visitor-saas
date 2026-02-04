FROM golang:1.25.6-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -v -o /run-app ./cmd/visitor

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/GeoLite2-Country.mmdb /usr/local/bin/
CMD ["run-app"]
