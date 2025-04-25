FROM golang:1.24 AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o freepdm ./cmd/pdmserver

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

EXPOSE 8080

COPY --from=builder /go/src/app/freepdm /app/freepdm

CMD ["/app/freepdm"]
