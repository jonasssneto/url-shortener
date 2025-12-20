FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/url-shortener ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/url-shortener /app/url-shortener

ENTRYPOINT ["/app/url-shortener", "-production"]