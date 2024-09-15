# CRUD Api пользователей

---
## Description
CRUD пользователей с авторизацией (access и refresh токены в cookie)
- GoLang 
  - Gin Framework
  - Gorm ORM
  - swaggo (Swagger Documentation 2.0)
  - golang-jwt (JSON Web Tokens)
- PostgreSQL
- Docker

## Get started
1. Запуск
    ```bash 
    docker-compose up --build
    ```
2. Переходим на [сайт api](http://localhost:8888/docs/index.html)
3. Создаем аккаунт (/registration)
4. Авторизация (/login) возвращает токены в cookie и body
5. Данные access токена (/acc-token)
6. Обновление пары access и refresh токенов (/refresh-token)

Для запуска не через докер 
1. Меняем файл `auth-service/.env`: `DB_HOST`, `DB_PORT`
2. БД
    ```bash
    docker-compose up --build db_auth 
    ```
3. Приложение
    ```bash
    cd auth-service
    ```
    ```bash
    go build main.go
    ```
    ```bash
    go run main.go
    ```