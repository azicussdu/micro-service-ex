# Зачем нужны Микросервисы:

- Независимые релизы — команда А обновила свой сервис и выкатила за 10 минут. Команда Б продолжает работать со своей версией.
- Устойчивость к отказам — если один сервис лёг, остальные могут работать (например, отключили рекомендации, но покупки идут).
- Разные технологии — сервис на Python для ML, на Rust для высоконагруженного прокси, на Java для корпоративной логики.
- Безопасная изоляция — компрометация сервиса истории заказов не даёт доступ к платежам.

# Пример микросервисов на Go

Этот проект содержит три Go-сервиса, которые запускаются совместно с помощью Docker Compose:

- `api-gateway`: API-шлюз на Gin с проверкой JWT и кешированием пользователей в Redis
- `user-service`: Сервис пользователей на Gin с PostgreSQL, GORM, регистрацией и входом в систему
- `order-service`: Сервис заказов на Gin с PostgreSQL, GORM и заказами в рамках пользователя

## Запуск

```bash
docker-compose up --build
```

Gateway(шлюз) доступен по адресу `http://localhost:8080`.

## Структура

```text
.
├── api-gateway
│   ├── cmd/api-gateway/main.go
│   ├── internal/config/config.go
│   ├── internal/handler/gateway.go
│   ├── internal/middleware/auth.go
│   ├── internal/service/user_cache.go
│   ├── internal/types/user.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── order-service
│   ├── cmd/order-service/main.go
│   ├── internal/config/config.go
│   ├── internal/handler/order_handler.go
│   ├── internal/middleware/internal_auth.go
│   ├── internal/model/order.go
│   ├── internal/repository/order_repository.go
│   ├── internal/service/order_service.go
│   ├── pkg/database/postgres.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── user-service
│   ├── cmd/user-service/main.go
│   ├── internal/config/config.go
│   ├── internal/handler/auth_handler.go
│   ├── internal/middleware/internal_auth.go
│   ├── internal/model/user.go
│   ├── internal/repository/user_repository.go
│   ├── internal/service/auth_service.go
│   ├── internal/service/token_service.go
│   ├── pkg/database/postgres.go
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── docker-compose.yml
└── README.md
```

## Маршруты шлюза

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/users/me`
- `POST /api/orders`
- `GET /api/orders`
- `DELETE /api/orders/:id`
- `GET /healthz`

## Маршруты сервиса пользователей

- `POST /auth/register`
- `POST /auth/login`
- `GET /users/:id`
- `GET /internal/users/:id`
- `GET /healthz`

## Маршруты сервиса заказов

- `POST /orders`
- `GET /orders`
- `DELETE /orders/:id`
- `GET /healthz`

## Примеры запросов

Регистрация:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'
```

Вход:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'
```

Используйте полученный токен:

```bash
export TOKEN="<jwt-from-login>"
```

Получение текущего пользователя:

```bash
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

Создание заказа:

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_name":"Mechanical Keyboard"}'
```

Список заказов:

```bash
curl http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN"
```

Удаление заказа:

```bash
curl -X DELETE http://localhost:8080/api/orders/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Примечания

- JWT генерируются сервисом пользователей и проверяются локально в шлюзе.
- Шлюз кеширует JSON пользователя в Redis после первого успешного обращения.
- Внутренние вызовы между сервисами используют X-Internal-Token.
- Сервисы пользователей и заказов выполняют AutoMigrate при запуске.
- Compose использует проверки здоровья (health checks) для Redis и обоих экземпляров PostgreSQL.
