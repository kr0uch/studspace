# StudSpace

HTTP API бэкенда VK Mini App: бронирование студенческих залов и модерация заявок администратором.

## Features

- Регистрация пользователя из VK launch params (`POST /auth/sync`); имя и фамилия только из VK
- Главные админы — `admins.txt`. При `POST /auth/sync` id из файла получает `admin`, остальные — `student`; убранный id теряет роль при следующем входе. Пустой файл, файл только с комментариями и отсутствие файла — админов нет. Некорректная строка не даёт серверу стартовать
- Профиль: отчество и должность (`PATCH /me`), флаг уведомлений (`PATCH /me/notifications`)
- Главная: свои активные брони за сегодняшние сутки МСК и залы со статусом `free` / `busy` / `closed`
- Бронь: `title` и `purpose` по 1–150 символов, `people_count` не больше вместимости зала. Живых заявок `pending` не больше 3
- Заявка `pending`, у которой `starts_at` уже наступил, в ответе становится `expired`: не входит в лимит, не одобряется. Одобренная бронь после `ends_at` в ответе — `completed`. Оба статуса не пишутся в БД
- Перенос одобренной брони создаёт новую заявку `pending` со ссылкой на исходную. Старый слот снимается только когда перенос одобрен. Пока первая заявка на перенос открыта, вторую создать нельзя
- Лента `GET /bookings` — только одобренные. Расписание показывает свои и чужие `pending` и `approved`. Занятость зала и конфликт слота считают только `approved`
- Модерация: список, approve, reject, cancel, return-to-review. После approve, reject и отмены админом владельцу уходит VK-уведомление, если флаг включён; сбой VK не отменяет операцию
- HMAC-подпись launch params, без сессий и JWT. CORS для Origin из конфига, `OPTIONS` отвечает 204 до проверки подписи
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
| `example.admins.txt` | образец списка главных админов. Скопируй в `admins.txt` |

`POSTGRES_PASSWORD` в `.env` и `postgres.postgres_password` в `.yaml` должны совпадать.

| Ключ | Описание |
|---|---|
| `vk.app_secret` | секрет приложения VK, HMAC launch params |
| `vk.notify_token` | сервисный токен `notifications.sendMessage`; пустой — не отправляем |
| `admins.txt` | VK ID главных админов, по одному на строку (не в `.yaml`). Пустые строки и `#` пропускаются. Пустой, только комментарии или отсутствующий файл — админов нет. Строка не из целого `≥ 1` — сервер не стартует. Файл в git не попадает |
| `server.cors_origins` | точные Origin фронта (`https://host`, со схемой и портом). Пусто — браузер с другого домена запросы не прочитает. Несколько через env: `CORS_ORIGINS=https://a,https://b` |
| `server.swagger.enabled` | `true` — открыть Swagger UI. По умолчанию выключен, в проде не включать. Env: `SWAGGER_ENABLED` |
| `server.swagger.host` | хост в Swagger UI для браузера (`localhost`), не `0.0.0.0` |
| `server.swagger.schema` | `http` / `https` |

## Usage

### Local

1. Подними PostgreSQL, скопируй `example.yaml` → `.yaml` и `example.admins.txt` → `admins.txt`. В `admins.txt` впиши VK ID главных админов.
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
| Swagger | http://localhost:8080/api/v1/swagger/ — редирект на `index.html`, только при `server.swagger.enabled: true` |
| Health | `GET /api/v1/health` (без подписи) |

### Docker

```bash
cp example.docker.yaml .yaml
cp example.env .env
cp example.admins.txt admins.txt
make docker-build
make docker-up
```

Снаружи слушает nginx на порту **80**. API внутри сети — `:8080`.

| | |
|---|---|
| API | http://localhost/api/v1/ |
| Swagger | http://localhost/api/v1/swagger/ — редирект на `index.html`, только при `server.swagger.enabled: true` |

Порядок сервисов: `postgres` → `migrate` → `studspace_server` → `nginx`. Конфиг монтируется как `.yaml` → `/app/.yaml`, список админов — `admins.txt` → `/app/admins.txt`, логи пишутся в `./logs`. Файл `admins.txt` должен существовать до `docker-up`, иначе Docker создаст вместо него каталог.

`nginx.conf` слушает HTTP **80** и проксирует `/api/v1/` на `studspace_server:8080` (30 r/s, burst 50). Остальное — 404. TLS для VK завершается перед контейнером: на хосте или у балансировщика. Внутри compose сертификата нет.

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

Подпись — HMAC-SHA256 строки `vk_*` как в query (`%2C`, `%20`), секрет `vk.app_secret`. `sign` в строку не входит. `vk_ts` старше 24 часов и `vk_user_id < 1` — `INVALID_SIGN`. Роль `admin` выставляет `POST /auth/sync` по `admins.txt`.

`GET /api/v1/health` и Swagger подпись не проверяют. Выключенный Swagger на `{prefix}/swagger/` отвечает 404, а не 401. Чужой `Origin` не получает `Access-Control-Allow-Origin`.

| Группа | Методы |
|---|---|
| Пользователь | `POST /auth/sync`, `GET /me`, `PATCH /me`, `PATCH /me/notifications` |
| Главная / залы | `GET /home`, `GET /spaces`, `GET /spaces/{id}/availability?date=` |
| Расписание | `GET /schedule`, `GET /calendar/busy-days` |
| Брони | `GET /bookings`, `POST /bookings`, `GET /bookings/my`, `POST /bookings/{id}/withdraw`, `POST /bookings/{id}/cancel`, `POST /bookings/{id}/reschedule` |
| Админ | `GET /admin/bookings`, `POST /admin/bookings/{id}/approve`, `.../reject`, `.../cancel`, `.../return-to-review` |

Полная схема — Swagger UI.
