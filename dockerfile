# Этап 1: Сборка приложения
FROM golang:1.21-alpine AS builder

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Этап 2: Запуск приложения
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Рабочая директория
WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Делаем файл исполняемым
RUN chmod +x ./main

# Переключаемся на непривилегированного пользователя
USER appuser

# Экспортируем порт
EXPOSE 8080

# Команда запуска
CMD ["./main"]