FROM golang:latest AS builder

WORKDIR /app

# Копируем go.mod и go.sum с проверкой существования
COPY go.mod go.sum* ./
RUN go mod download || true

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]