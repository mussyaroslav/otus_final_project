# Используем официальный образ Go
FROM golang:1.14-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod download

# Собираем бинарник
RUN go build -o main .

# Запуск приложения
CMD ["./main"]
