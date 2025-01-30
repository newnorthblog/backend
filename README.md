# New North Backend

## Настройка

### Предварительные требования

- Go 1.18+
- Docker
- Docker Compose

### Переменные окружения

Создайте файл `.env` в корневом каталоге и добавьте следующие переменные:

```
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_HOST=your_postgres_host
POSTGRES_PORT=your_postgres_port
POSTGRES_DB=your_postgres_db
MIGRATION_DIR=your_migration_directory
```

## Команды Makefile

### Установка зависимостей

Установите необходимые инструменты и зависимости.

```sh
make install-goose
make install-golangci-lint
```

### Линтер

Запустите линтер.

```sh
make lint
```

### Docker Compose

Запустите и остановите службы Docker Compose.

```sh
make compose
make compose-down
```

### Сборка

Соберите бинарный файл backend.

```sh
make build
```

### Запуск

Запустите приложение.

```sh
make run
```

### Swagger

Сгенерируйте документацию Swagger.

```sh
make swag
```

### Миграции

Управляйте миграциями базы данных.

```sh
make migration-status
make migration-up
make migration-down
make migration-create
```