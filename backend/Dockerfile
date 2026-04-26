FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Находим main.go и собираем
RUN go build -o main $(find . -name "main.go" | head -1)

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]