# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Финальный этап
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из сборщика
COPY --from=builder /app/main .
# Копируем статические файлы (если есть)
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]