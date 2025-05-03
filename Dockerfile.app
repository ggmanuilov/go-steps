FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o delivery-server internal/server.go

FROM alpine:3.21
WORKDIR /opt/
COPY --from=builder /app/delivery-server .
COPY --from=builder /app/.env.example .env
RUN apk --no-cache add ca-certificates
EXPOSE 8080
ENV ENV_FILE=/opt/.env
CMD ["./delivery-server"]