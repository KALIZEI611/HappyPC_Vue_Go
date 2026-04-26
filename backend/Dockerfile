FROM golang:latest AS builder

WORKDIR /app

# Проверяем наличие файлов и копируем
COPY go.mod go.sum* ./
RUN if [ -f go.mod ]; then go mod download; else echo "go.mod not found"; exit 1; fi

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]