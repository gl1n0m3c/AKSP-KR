FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/migrator ./services/migrator/cmd/migrator

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/migrator /app/migrator
COPY services/migrator/config/config.yaml /app/config.yaml

ENV CONFIG_PATH=/app/config.yaml
CMD ["/app/migrator"]

