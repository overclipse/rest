# Используем Go 1.25
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Кладём только mod/sum для кеша зависимостей
COPY auth/go.mod auth/go.sum ./auth/

# Если в go.mod есть replace на локальный jwtmanager — нужен каталог
#   replace jwt_manager => ../pkg/jwtmanager
COPY pkg ./pkg

WORKDIR /app/auth
RUN go mod download

# Копируем исходники
COPY auth/ .

# Собираем
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth .

# Рантайм
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/bin/auth ./auth
EXPOSE 8101
CMD ["./auth"]