# Тестовое задание на стажировку в Авито (Backend)

Сервис назначения ревьюеров для Pull Request’ов. Написан на языке Go, в качестве базы данных используется PostgreSQL.

## Сборка и запуск

Для сборки и запуска приложения введите следующую команду:

```bash
make up
```

Сервис будет запущен по адресу localhost:8080

Для просмотра логов приложения введите следующую команду:

```bash
make logs
```

Чтобы провести интеграционные тесты введите:

```bash
make tests
```

Остановить сервис можно с помощью команды:

```bash
make down
```

## Используемые инструменты

* [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - программа для генерации кода сервера и клиента на основе OpenAPI спецификации. Так же автоматически генерирует код для парсинга запросов и сериализации ответов
* [Gin](https://github.com/gin-gonic/gin) - фреймворк для написания веб-приложений
* [pgx](https://github.com/jackc/pgx) - драйвер для работы с PostgreSQL
* [go-transaction-manager](https://github.com/avito-tech/go-transaction-manager) - менеджер транзакций
* [cleanenv](https://github.com/ilyakaznacheev/cleanenv) - библиотека для чтения файлов конфигурации
* [migrate](https://github.com/golang-migrate/migrate) - программа и библиотека для управления миграциями БД
* [testify](https://github.com/stretchr/testify) - библиотека для написания тестов
* [gofakeit](https://github.com/brianvoe/gofakeit) - библиотека для генерации случайной информации, как в форме чисес, так и в форме текста (случайные ID, слова, предложения)