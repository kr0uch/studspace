# StudSpace

HTTP API бэкенда VK Mini App: бронирование студенческих залов, запись на мероприятия и модерация заявок администратором.

## Features

- Регистрация пользователя из VK launch params (`POST /auth/sync`)
- Брони залов: создание, список, отзыв, отмена, перенос, лимит pending
- Модерация заявок админом (approve / reject / cancel / return-to-review)
- Мероприятия: лента, запись, лист ожидания, уведомления VK
- Расписание и занятость залов
- HMAC-подпись launch params, без сессий и JWT
- Swagger UI, миграции PostgreSQL, Docker Compose + nginx

## Tech stack

| | |
|---|---|
| Language | Go 1.27 |
| HTTP | `net/http`, swaggo |
| DB | PostgreSQL 17, sqlx, pgx v5, golang-migrate |
| Log | zap |
| Proxy | nginx |
| Tests | testify, sqlmock |

## Requirements

- Go 1.27+
- PostgreSQL 17 (локальный запуск)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- Docker и Docker Compose (запуск в контейнерах)
- Make
- [swag](https://github.com/swaggo/swag) — только для перегенерации доки

## Installation

```bash
git clone https://github.com/kr0uch/studspace.git
cd studspace
```

Шаблоны конфигов:

```bash
cp example.yaml .yaml                 # локально
# или
cp example.docker.yaml .yaml          # Docker
cp example.env .env
```

`.yaml` и `.env` в git не попадают. Заполни секреты (см. [Configuration](#configuration)).

## Configuration

Приложение читает `.yaml` из корня (`internal/config`). Compose дополнительно берёт `.env`.

| Файл | Назначение |
|---|---|
| `example.yaml` | локальный запуск: `postgres_host: localhost`, `server.host: localhost` |
| `example.docker.yaml` | контейнер: `postgres_host: postgres`, `server.host: 0.0.0.0` |
| `example.env` | переменные Compose: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_SSLMODE` |

`POSTGRES_PASSWORD` в `.env` и `postgres.postgres_password` в `.yaml` должны совпадать.

| Ключ | Описание |
|---|---|
| `vk.app_secret` | секрет приложения VK, HMAC launch params |
| `vk.notify_token` | сервисный токен `notifications.sendMessage`; пустой — не отправляем |
| `server.swagger.host` | хост в Swagger UI для браузера (`localhost`), не `0.0.0.0` |
| `server.swagger.schema` | `http` / `https` |

## Usage

### Local

1. Подними PostgreSQL и скопируй `example.yaml` → `.yaml`.
2. Миграции:

```bash
migrate -path migrations -database "postgres://postgres:password@localhost:5432/studspace?sslmode=disable" up
```

3. Запуск:

```bash
go run ./cmd/main.go
```

| | |
|---|---|
| API | http://localhost:8080/api/v1/ |
| Swagger | http://localhost:8080/api/v1/swagger/index.html |
| Health | `GET /api/v1/health` (без подписи) |

### Docker

```bash
cp example.docker.yaml .yaml
cp example.env .env
make docker-build
make docker-up
```

Снаружи слушает nginx на порту **80**. API внутри сети — `:8080`.

| | |
|---|---|
| API | http://localhost/api/v1/ |
| Swagger | http://localhost/api/v1/swagger/index.html |

Порядок сервисов: `postgres` → `migrate` → `studspace_server` → `nginx`. Конфиг монтируется как `.yaml` → `/app/.yaml`, логи пишутся в `./logs`.

`nginx.conf` проксирует `/api/v1/` на `studspace_server:8080` (30 r/s, burst 50). Остальное — 404.

Остановка: `make docker-down`. С удалением томов: `make docker-clean`.

## Development

```bash
make test
make swag          # docs/docs.go, swagger.json, swagger.yaml
make docker-build
make docker-up
make docker-logs
make docker-down
make docker-restart
make docker-clean  # down -v + prune, удаляет данные postgres
```

## API

Базовый путь: `/api/v1`.

Все ручки, кроме health и swagger, требуют заголовок:

```http
Authorization: VK vk_user_id=...&vk_ts=...&sign=...
```

Подпись — HMAC-SHA256 от launch params (`vk.app_secret`). `vk_user_id` только из params. Роль `admin` — колонка в БД.

| Группа | Методы |
|---|---|
| Пользователь | `POST /auth/sync`, `GET /me`, `PATCH /me/notifications` |
| Главная / залы | `GET /home`, `GET /spaces`, `GET /spaces/{id}/availability` |
| Расписание | `GET /schedule`, `GET /calendar/busy-days` |
| Брони | `POST /bookings`, `GET /bookings/my`, withdraw / cancel / reschedule |
| Мероприятия | `GET /events`, registration, waitlist |
| Админ | approve / reject / cancel / return-to-review, CRUD мероприятий |

Полная схема — Swagger UI.
