FROM golang:1.22-alpine AS builder

WORKDIR /my-app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/goods ./cmd/goods

FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/goods .

# Копируем конфигурационный файл
COPY ./config/local.yaml ./config/local.yaml

# Устанавливаем команду по умолчанию для запуска приложения
CMD ["./goods", "--config=./config/local.yaml"]