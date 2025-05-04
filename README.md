# Todo List Telegram Bot

Telegram бот для управления списком задач с возможностью добавления, удаления, изменения и просмотра задач.

## Возможности

- Добавление новых задач
- Удаление существующих задач
- Изменение статуса и описания задач
- Просмотр списка всех задач
- Категоризация задач
- Установка дедлайнов

## Технологии

- Go 1.24
- PostgreSQL
- Telegram Bot Token

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/dunooo0ooo/todo-tg-bot.git
cd todo-tg-bot
```

2. Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

3. Заполните необходимые переменные окружения в файле `.env`:
```
BOT_TOKEN=your_telegram_bot_token
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todo_db
```

4. Установите зависимости:
```bash
go mod download
```

5. Запустите миграции базы данных:
```bash
go run ./cmd/migrate/main.go
```

6. Запустите бота:
```bash
docker compose up -d
```

## Использование

После запуска бота, вы можете использовать следующие команды:

- `/start` - Начать работу с ботом
- `/help` - Показать список доступных команд
- `/add_task` - Добавить новую задачу
- `/delete_task` - Удалить задачу
- `/change_task` - Изменить существующую задачу
- `/show_task` - Показать список всех задач


### Запуск тестов
```bash
go test -v ./...
```

### TODO
- Add notification service using Kafka