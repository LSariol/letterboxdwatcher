FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o letterboxdwatcher ./cmd/letterboxdwatcher

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/letterboxdwatcher .
CMD ["./letterboxdwatcher"]
