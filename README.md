# Запуск🚀
## Перед запуском приложения, создайте в корневой директории файл .env (с заменой моего), и заполните поля:
```
PGUSER=
PGPASSWORD=
PGHOST=localhost
PGPORT=5436
PGDATABASE=
PGSSLMODE=disable
HTTP_PORT=8082
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=43200m # 1 month
SIGNING_KEY=qazwsxedc
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=
SMTP_PASS=
```

## Запустите PostgreSQL в отдельном Docker-контейнере с указанием ваших настроек
Например:
```
docker run --name medods-test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=medods-test -p 5436:5432 -d postgres
```
Установите зависимости для приложения:
go mod download

Также перед запуском нужно выполнить миграцию для таблицы с заметками. Можно использовать эту [библиотеку](https://github.com/golang-migrate/migrate)
Пример команды:
```
migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5436/medods-test?sslmode=disable' up
```
