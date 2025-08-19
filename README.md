# Тестовое задание  
## Демонстрационный сервис с Kafka, PostgreSQL, кешем


## Запуск:
1. Клонировать репозиторий
 ```bash
    git clone https://github.com/AvenidAcami/wb_lvl0.git
    cd wb_lvl0
```
2. Создать .env в соответствии с .env.example
3. Запустить с помощью
```bash
    docker compose up --build
```

## Важно:
 - Если нужно просто протестировать, то в .env в DB_USER нужно вписать wb_lvl0 и в DB_PASSWORD вписать test
   (либо менять миграцию на создание пользователя и вставлять свои данные)

## Использование:
 - Получения заказа: 
```http request
    GET /order/<order_uid>
```

 - Документация:
```http request
    GET /swagger/index.html
```

## Используемые технологии:
 - Golang
 - Docker
 - Docker compose
 - Redis

## Структура
 - `cmd/app/` — точка входа в приложение
 - `config/` — инициализация сторонних компонентов
   - `migrations/` — скрипты для миграций
 - `docs/` — swagger документация
 - `internal/` — основная внутренняя логика проекта:
   - `app/` — инициализация основных компонентов
   - `controller/` — принятие запросов, вызов сервиса
   - `kafka/` — консьюмер кафки
   - `model/` — структуры данных
   - `repository/` — работа с бд
   - `router/` — определение маршрутов
   - `service/` — бизнес-логика
 - `tests/` — тесты
   - `integration/` — интеграционные тесты
   - `unit/` — unit-тесты
 - `web/` — веб-интерфейс
