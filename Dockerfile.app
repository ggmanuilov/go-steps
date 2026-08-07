FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
# Чистый Go-резолвер: без него cgo-резолвер (debian) берёт IPv6 первым
# и падает с "network is unreachable", т.к. в Docker-сети нет IPv6-маршрута.
ENV GODEBUG=netdns=go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o delivery-server internal/server.go

FROM alpine:3.21
WORKDIR /opt/
COPY --from=builder /app/delivery-server .
COPY --from=builder /app/.env.example .env
RUN apk --no-cache add ca-certificates
EXPOSE 8080
ENV ENV_FILE=/opt/.env
CMD ["./delivery-server"]