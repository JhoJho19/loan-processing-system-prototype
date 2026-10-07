# Loan Processing System Prototype


## English

A Go pet project: a backend prototype for accepting loan applications. The service validates customer and application data, assigns IDs and the `NEW` status, then stores both records in PostgreSQL in a single transaction. Looking up an application by ID is also implemented.

### Tests

- Service unit tests cover application creation, invalid input, repository errors, and application lookup (including an empty ID).
- PostgreSQL integration tests cover saving a customer and application, transaction rollback on failure, duplicate customer ID conflicts, lookup by ID, and not-found results.
- Integration tests start PostgreSQL 16 in Docker via Testcontainers. Run all tests with `go test ./...`; skip integration tests with `go test -short ./...`.


## Русский

Пет-проект на Go: прототип бэкенда для приема заявок на кредит. Сервис проверяет данные клиента и заявки, присваивает им идентификаторы и статус `NEW`, а затем сохраняет их в PostgreSQL в одной транзакции. Также реализовано получение заявки по ID.

### Тесты

- Юнит-тесты сервиса: создание заявки, отклонение некорректного ввода, обработка ошибок репозитория и получение заявки (включая пустой ID).
- Интеграционные тесты PostgreSQL: сохранение клиента и заявки, откат транзакции при ошибке, конфликт дублирующегося ID клиента, получение заявки по ID и сценарий «не найдено».
- Интеграционные тесты запускают PostgreSQL 16 в Docker через Testcontainers. Запуск всех тестов: `go test ./...`; без интеграционных тестов: `go test -short ./...`.

