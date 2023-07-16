FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app/server .


FROM alpine:3.18
WORKDIR /opt/delivery/
COPY --from=builder ./app/ .
COPY --from=builder ./app/.env.example .env
EXPOSE 8080
CMD ["./app/server"]