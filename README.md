# TGNotification Service

Микросервис для асинхронной отправки уведомлений пользователям в Telegram. Часть системы маркетплейса "Зоомагазин".

## Технологический стек
- **Language:** Go 1.21+
- **Architecture:** Clean Architecture (Layers: Transport, Service, Repository)
- **Communication:**
    - **Async:** Redis (Pub/Sub & Queues) для отправки уведомлений.
    - **Sync:** gRPC (Protocol Buffers) для синхронизации аккаунтов пользователей.
- **Observability:** Structured logging (`slog`) с использованием контекста.

## Структура проекта
- `cmd/`: Точка входа приложения.
- `internal/`: Бизнес-логика (Service), работа с БД (Repository), транспорт (Telegram, gRPC).
- `pkg/`: Общие утилиты (Config, Logger, API).
- `proto/`: Контракты gRPC.

## Как запустить
1. Установите зависимости: `go mod download`.
2. Настройте `.env` файл (пример в `.env.example`).
3. Запустите Redis и gRPC-сервер Django.
4. Запуск: `go run cmd/bot/main.go`.