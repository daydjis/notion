# Stage 1: Build Go app
FROM golang:1.22.5-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o server ./cmd/main.go

# Stage 2: Runtime + migrate
FROM alpine:latest

# Установка утилиты migrate
RUN apk add --no-cache ca-certificates curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    | tar xz && mv migrate /usr/bin/migrate

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .

EXPOSE 8080

# Выполняем миграции и запускаем сервер
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]
