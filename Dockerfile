FROM golang:1.23 AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o freepdm ./cmd/pdmserver/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /go/src/app/freepdm /app/freepdm

CMD ["/app/freepdm"]
