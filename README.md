# Тестовое задание  
## Демонстрационный сервис с Kafka, PostgreSQL, кешем


## Запуск:
1. Клонировать репозиторий
 ```bash
    git clone github.com/AvenidAcami/wb_lvl0.git
    cd wb_lvl0
```
2. Создать .env в соответствии с .env.example
3. Запустить с помощью
```bash
    docker compose up --build
```

## Использование:
 - Получения заказа: 
```http request
    GET /order/<order_uid>
```

 - Документация:
```http request
    GET /swagger/index.html
```

## Важно:
 - Если нужно просто протестировать, то в .env в DB_USER нужно вписать wb_lvl0 м в DB_PASSWORD вписать test
   (либо менять миграцию на создание пользователя и вставлять свои данные)
