FROM golang:1.23.3 AS builder

WORKDIR /review

# Устанавливаем pg_isready и очищаем кеш
RUN apt-get update && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

# Отключаем CGO для статической компиляции
 ENV CGO_ENABLED=0

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /review/review_service ./cmd/server.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /review

# Устанавливаем postgresql-client и dos2unix
RUN apk add --no-cache postgresql-client dos2unix

COPY .env /review/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /review/review_service /review/review_service

# Копируем wait-for-db.sh и делаем исполняемым
COPY --from=builder /review/wait-for-db.sh /review/wait-for-db.sh
RUN chmod +x /review/wait-for-db.sh

# Преобразуем формат строки в скрипте wait-for-db.sh в Unix-формат
RUN dos2unix /review/wait-for-db.sh

# Запуск приложения
CMD ["/review/review_service"]
