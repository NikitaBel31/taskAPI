# TaskAPI

Простой REST API для управления задачами на Go.  
Поддерживает создание задач, получение по ID и список с фильтрацией по статусу.  
Логирование реализовано через асинхронный самописный логгер, есть middleware для Request ID и таймаута запросов.   

## Архитектура
Проект построен по принципу чистой архитектуры:
- **cmd/** — точка входа в приложение.
- **internal/app** — инициализация зависимостей.
- **internal/config** — конфигурация сервиса.
- **internal/domain** — доменные сущности.
- **internal/dto** — структуры запросов/ответов.
- **internal/handlers/http** — HTTP-обработчики.
- **internal/logger** — асинхронный JSON-логгер.
- **internal/repository** — интерфейсы репозиториев и их реализации.
- **internal/usecase** — бизнес-логика.
- **tests** — Unit-тесты.

## Возможности
- **Создание задачи** (`POST /tasks`)
- **Получение задачи по ID** (`GET /tasks/{id}`)
- **Список задач** (`GET /tasks`)
- **Список задач с фильтрацией по статусу** (`GET /tasks?status={status}`)
- **Проверка работоспособности** (`GET /health`)

## Статусы задач
- `todo`
- `in_progress`
- `done`

## Запуск
```bash
git clone https://github.com/NikitaBel31/taskAPI.git
cd taskAPI
go run ./cmd/task-service/main.go
```
## Примеры запросов
### 1. Создание задачи
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test task",
    "description": "test description",
    "status": "in_progress"
  }'
```
### 2. Получение списка задач
```curl -X GET http://localhost:8080/tasks```

### 3. Получение списка задач с фильтрацией по статусу
```curl -X GET http://localhost:8080/tasks?status=in_progress```

### 4. Получение задачи по ID
```curl -X GET http://localhost:8080/tasks/{id}```

### 5. Проверка работоспособности
```curl -X GET http://localhost:8080/health```
