FROM golang:latest AS builder

WORKDIR /app

# Копируем всё содержимое backend
COPY . .

# Скачиваем зависимости (теперь go.mod есть)
RUN go mod download

# Собираем бинарник, указывая правильный путь
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]