# 1. Используем официальный базовый образ Go
FROM golang:1.24 as builder

# 2. Создаем рабочую папку внутри контейнера
WORKDIR /app

# 3. Копируем модули зависимостей
COPY go.mod go.sum ./

# 4. Устанавливаем зависимости
RUN go mod download

# 5. Копируем исходный код проекта
COPY . .
COPY .env .env


# 6. Собираем приложение в бинарный файл
RUN go build -o bookingkart-platform ./cmd/main.go

# 7. Минимизируем финальный образ
FROM debian:bookworm-slim

# 8. Создаем рабочую папку в финальном контейнере
WORKDIR /app

# 9. Копируем готовый бинарник из builder-стадии
COPY --from=builder /app/bookingkart-platform .

# 10. Открываем нужный порт
EXPOSE 8080

# 11. Запускаем приложение
ENTRYPOINT ["./bookingkart-platform"]
