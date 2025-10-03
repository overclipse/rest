# Используем официальный образ Go как базовый
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go mod файлы auth сервиса
COPY auth/go.mod auth/go.sum ./auth/

# Копируем pkg директорию для локального модуля
COPY pkg ./pkg

# Переходим в папку auth
WORKDIR /app/auth

# Загружаем зависимости
RUN go mod download

# Копируем исходный код auth сервиса
COPY auth/ .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный этап - минимальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл из builder этапа
COPY --from=builder /app/auth/main .

# Открываем порт
EXPOSE 8101

# Запускаем приложение
CMD ["./main"]