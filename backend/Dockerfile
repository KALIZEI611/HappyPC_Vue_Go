FROM golang:latest AS builder

WORKDIR /app

# Копируем go.mod и go.sum (они лежат рядом с Dockerfile)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код бэкенда
COPY . .

# Собираем бинарник, указывая путь к main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Финальный этап
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из сборщика
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]