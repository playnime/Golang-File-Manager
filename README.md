# Golang File Manager

Защищённый файловый менеджер на Go: REST API для хранения, скачивания и обмена файлами между пользователями — с шифрованием на клиенте, обменом ключами по RSA и двухфакторной аутентификацией (TOTP).

> Учебный проект: криптографические примитивы (RSA, тест Миллера — Рабина, MD5) реализованы с нуля вручную, без использования готовых библиотек.

## Возможности

- Регистрация и вход по логину/паролю, сессии с Bearer-токенами (время жизни — 5 минут)
- Двухфакторная аутентификация через TOTP (Google Authenticator и совместимые приложения)
- Обмен ключами: клиент отправляет свой RSA-открытый ключ, сервер возвращает AES-128 ключ, зашифрованный этим ключом
- Файлы передаются по сети в зашифрованном виде (AES-GCM), целостность проверяется собственным MD5-хешем, добавляемым к данным
- Загрузка, получение списка, скачивание и удаление файлов
- Передача прав на файлы другим пользователям (`gvRights`)
- PostgreSQL + SQL-миграции, Docker Compose для запуска всего стека

## Архитектура

```
┌────────┐   1. POST /keys/exchange (N, E)      ┌─────────────────────────────┐
│ Клиент │ ───────────────────────────────────► │        Сервер (Go)          │
│        │   2. AES-ключ = D(C^d mod n)         │                             │
│        │ ◄─────────────────────────────────── │  генерирует AES-128 ключ,   │
│        │      C = AES^e mod n (hex)           │  шифрует его через RSA,     │
│        │                                      │  сохраняет за пользователем │
│        │   3. upload: AES-GCM(data + MD5)     │                             │
│        │ ───────────────────────────────────► │  расшифровывает, проверяет  │
│        │   4. download: AES-GCM(data + MD5)   │  хеш, хранит файлы на диске │
│        │ ◄─────────────────────────────────── │  и метаданные в PostgreSQL  │
└────────┘                                      └─────────────────────────────┘
```

Сервер собран по чистой архитектуре: `transport → service → repository`.

```
server/
├── cmd/server/main.go            # точка входа
└── internal/
    ├── core/
    │   ├── httpServer/            # роутер, middleware авторизации
    │   ├── Aes/                   # AES-GCM + проверка целостности через MD5
    │   ├── RSA/                   # собственная RSA: генерация, шифрование, Miller-Rabin
    │   ├── md5MY/                 # собственная реализация MD5
    │   └── db/                    # подключение к БД, миграции
    ├── users/                     # регистрация, логин, сессии, 2FA, обмен ключами
    └── files/                     # загрузка/скачивание/удаление, права доступа
client/main.go                     # пример клиента: заливка файла на /add
migrations/                        # SQL-миграции (golang-migrate)
```

## API

### Публичные эндпоинты

| Метод | Путь | Тело | Ответ |
|---|---|---|---|
| POST | `/register` | `{"login": "...", "password": "..."}` | `201` + логин |
| POST | `/login` | `{"login": "...", "password": "..."}` | `{"token": "..."}` или `{"tmp_token": "...", "require_2fa": true}` |
| POST | `/login/2fa` | `{"tmp_token": "...", "code": "123456"}` | `{"token": "..."}` |

Логин и пароль — от 4 до 20 символов.

### Защищённые эндпоинты (заголовок `Authorization: Bearer <token>`)

| Метод | Путь | Тело | Ответ |
|---|---|---|---|
| POST | `/keys/exchange` | `{"N": "<hex>", "E": "<hex>"}` | `{"EncryptedAESKey": "<hex>"}` |
| POST | `/add` | multipart/form-data, поле `File` (AES-GCM + MD5) | `201` |
| GET | `/get` | — | JSON со списком файлов |
| GET | `/getfile` | `{"Filename": "...", "Author": "..."}` | multipart с зашифрованным файлом |
| DELETE | `/delete` | `{"Filename": "...", "Author": "..."}` | `204` |
| POST | `/gvRights` | `{"Filename": "...", "User": "...", "Rights": 1}` | текст |
| GET | `/2fa/setup` | — | `otpauth://` URL + секрет |
| POST | `/2fa/confirm` | `{"code": "123456"}` | 2FA включена |
| DELETE | `/2fa` | `{"code": "123456"}` | 2FA отключена |

## Запуск

### Через Docker (рекомендуется)

Нужны Docker и `make`.

1. Создайте файл `.env` в корне проекта:

    ```env
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=lab2pepe
    CONN_STR=postgres://postgres:postgres@db:5432/lab2pepe?sslmode=disable
    ```

2. Поднимите стек:

    ```bash
    make up        # postgres + миграции + сервер
    make logs      # логи сервера
    make down      # остановить
    make clean     # остановить и удалить volumes (включая файлы и БД)
    ```

Сервер слушает порт `9091`. Файлы хранятся в volume `uploads` по пути `$UPLOAD_DIR/<login>/<filename>`.

### Локально

```bash
make migrate-up    # применить миграции (нужен бинарник migrate и CONN_STR в .env)
make service-run   # go run server/cmd/server/main.go
```

### Пример клиента

`client/main.go` — минимальный клиент: отправляет один файл на `POST /add` с Bearer-токеном (файл должен лежать рядом с `main.go`). Токен получается при регистрации/логине.

## Стек

- Go 1.26, `net/http` + [gorilla/mux](https://github.com/gorilla/mux)
- PostgreSQL + [pgx/v5](https://github.com/jackc/pgx), миграции — [golang-migrate](https://github.com/golang-migrate/migrate)
- TOTP — [pquerna/otp](https://github.com/pquerna/otp)
- Криптография: AES-GCM из стандартной библиотеки; RSA, тест Миллера — Рабина и MD5 — собственная реализация

## Лицензия

[MIT](LICENSE)
