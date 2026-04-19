# backend/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum из папки backend (контекст сборки - корень репозитория)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Копируем весь код бэкенда
COPY backend/ .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Финальный этап
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из сборщика
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]