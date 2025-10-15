FROM golang:1.20-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -o /photobooth-api ./cmd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /photobooth-api /photobooth-api

EXPOSE 8080
ENTRYPOINT ["/photobooth-api"]
