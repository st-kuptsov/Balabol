# ===========================
# Сборка приложения
# ===========================
FROM golang:1.24-alpine AS builder

# Устанавливаем базовые инструменты
RUN apk add --no-cache git ca-certificates

WORKDIR /usr/src/mail2tg

# Копируем модули и исходники
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY config ./config
COPY internal ./internal
COPY pkg ./pkg

# Версия приложения
ARG APP_VERSION=dev

# Собираем бинарник статически для минимального образа
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X main.Version=${APP_VERSION}" \
    -o ./bin/balabol ./cmd/balabol/main.go

# ===========================
# Минимальный runtime контейнер с поддержкой таймзоны
# ===========================
FROM alpine:latest

WORKDIR /app

# Устанавливаем curl для healthcheck и tzdata + timedatectl
RUN apk add --no-cache curl tzdata

# Копируем бинарник из builder
COPY --from=builder /usr/src/mail2tg/bin/balabol ./balabol

# Копируем конфигурацию и secrets
COPY config/config.example.yaml ./config/config.yaml
COPY config/secrets.example.yaml ./config/secrets.yaml

# Настройка таймзоны через переменную окружения TZ
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Прокидываем порт для метрик
EXPOSE 9090/tcp

# Запуск приложения
CMD ["/app/balabol", "-config", "/app/config/config.yaml"]
