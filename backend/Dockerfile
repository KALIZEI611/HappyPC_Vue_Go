FROM golang:latest AS builder

WORKDIR /app

# Копируем go.mod и go.sum из папки backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Копируем остальной код из папки backend
COPY backend/ .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]