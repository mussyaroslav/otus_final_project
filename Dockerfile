# Используем официальный образ Golang
FROM golang:1.22.2 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Сборка приложения
RUN go build -o main .

# Используем легкий образ для запуска
FROM alpine:latest

# Копируем скомпилированный бинарник из builder
COPY --from=builder /app/main .

# Устанавливаем необходимые зависимости
RUN apk add --no-cache ca-certificates

# Запускаем приложение
ENTRYPOINT ["./main"]
