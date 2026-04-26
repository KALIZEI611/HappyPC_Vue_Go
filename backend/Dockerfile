FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Сборка сама скачает зависимости
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]