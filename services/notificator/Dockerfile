FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/notificator ./services/notificator/cmd/notificator

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/notificator /app/notificator
COPY services/notificator/config/config.yaml /app/config.yaml

ENV CONFIG_PATH=/app/config.yaml
CMD ["/app/notificator"]

