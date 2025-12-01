# История деплоя Savory AI Server на Railway

## Что было сделано

### 1. Исправление CORS для QR-кодов
- Обновлён `app/module/qr_code/controller/qr_code_controller.go`
- Заменён `ctx.SendFile()` на `os.ReadFile()` + `ctx.Send(data)`
- Добавлен явный заголовок `Access-Control-Allow-Origin: *`
- Обновлён CORS middleware в `internal/bootstrap/webserver.go`

### 2. Подготовка к Railway

#### Созданные файлы:
- `railway.toml` - конфигурация Railway
- `config/config.railway.toml` - production конфиг с переменными окружения

#### Изменённые файлы:
- `utils/config/config.go` - добавлена поддержка переменных окружения в формате `${VAR}` и `${VAR:-default}`
- `Dockerfile` - обновлён до Go 1.24, добавлено копирование railway конфига
- `config/config.docker.toml` - API ключ заменён на `${ANTHROPIC_API_KEY}`
- `.gitignore` - исправлено для разрешения docker и railway конфигов

### 3. Деплой на Railway

#### Команды:
```bash
railway init          # Инициализация проекта
railway add           # Добавление PostgreSQL
railway service       # Привязка сервиса
railway up            # Деплой
railway domain        # Получение публичного домена
```

#### Необходимые переменные окружения в Railway:
```
CONFIG_FILE=/app/config/config.railway.toml
JWT_SECRET=твой-jwt-секрет
ANTHROPIC_API_KEY=sk-ant-api03-...
DATABASE_URL=автоматически от PostgreSQL
PORT=автоматически от Railway
```

### 4. Структура config.railway.toml
```toml
[app]
name = "Savory AI Server"
host = "0.0.0.0"
port = ":${PORT:-4000}"
production = true
CHAT_SERVICE_URL = "https://${RAILWAY_PUBLIC_DOMAIN}"

[db.postgres]
dsn = "${DATABASE_URL}"

[middleware.jwt]
secret = "${JWT_SECRET}"

[anthropic]
api_key = "${ANTHROPIC_API_KEY}"
```

### 5. Проблемы и решения

| Проблема | Решение |
|----------|---------|
| CORS ошибка для QR-кодов | Явный `Access-Control-Allow-Origin: *` + `ctx.Send()` вместо `SendFile()` |
| `*.toml` в .gitignore | Исправлен gitignore для разрешения docker/railway конфигов |
| GitHub блокирует пуш с API ключом | Заменён ключ на `${ANTHROPIC_API_KEY}` |
| Go 1.24 не найден в Dockerfile | Обновлён `FROM golang:1.23-alpine` на `golang:1.24-alpine` |


## Важно

- Локальный `config/config.toml` НЕ коммитится (содержит секреты)
- API ключи передаются ТОЛЬКО через переменные окружения
- PostgreSQL подключается автоматически через `DATABASE_URL`
