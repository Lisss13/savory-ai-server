# CLAUDE.md - Savory AI Server

## Обзор проекта

Go REST API сервер для системы управления ресторанами. Построен на Clean Architecture с использованием Uber FX для dependency injection.

**Основной стек:** Go, Fiber v2, GORM, PostgreSQL, JWT

## Команды

```bash
# Запуск с миграциями
go run cmd/main.go -migrate

# Запуск с сидированием БД
go run cmd/main.go -seed

# Docker
docker-compose up -d

# Установка зависимостей
go mod tidy
```

## Архитектура

### Слои приложения

```
Controller (HTTP) → Service (бизнес-логика) → Repository (БД)
```

### Структура директорий

```
server/
├── cmd/main.go                  # Точка входа, регистрация FX модулей
├── app/
│   ├── middleware/              # HTTP middleware
│   ├── module/                  # Бизнес-модули
│   ├── router/api.go            # Регистрация маршрутов
│   └── storage/                 # GORM модели
├── internal/bootstrap/          # Инициализация приложения
│   └── database/                # Подключение к PostgreSQL
├── utils/
│   ├── config/                  # Конфигурация (TOML)
│   ├── jwt/                     # JWT утилиты
│   └── response/                # Форматирование ответов
└── config/                      # TOML конфиги
```

## Создание нового модуля

### 1. Создать структуру директорий

```
app/module/{module_name}/
├── {module_name}_module.go
├── controller/
│   ├── controller.go
│   └── {module_name}_controller.go
├── service/
│   └── {module_name}_service.go
├── repository/
│   └── {module_name}_repository.go
└── payload/
    ├── request.go
    └── response.go
```

### 2. Модель в storage (app/storage/{module_name}.go)

```go
package storage

import "gorm.io/gorm"

type Example struct {
    gorm.Model
    Name        string `gorm:"column:name;not null" json:"name"`
    Description string `gorm:"column:description" json:"description"`
}
```

### 3. Repository (app/module/{name}/repository/{name}_repository.go)

```go
package repository

import (
    "savory-ai-server/app/storage"
    "savory-ai-server/internal/bootstrap/database"
)

type exampleRepository struct {
    DB *database.Database
}

type ExampleRepository interface {
    FindAll() ([]storage.Example, error)
    FindByID(id uint) (*storage.Example, error)
    Create(example *storage.Example) (*storage.Example, error)
    Update(example *storage.Example) (*storage.Example, error)
    Delete(id uint) error
}

func NewExampleRepository(db *database.Database) ExampleRepository {
    return &exampleRepository{DB: db}
}

func (r *exampleRepository) FindAll() ([]storage.Example, error) {
    var examples []storage.Example
    err := r.DB.DB.Find(&examples).Error
    return examples, err
}

func (r *exampleRepository) FindByID(id uint) (*storage.Example, error) {
    var example storage.Example
    err := r.DB.DB.First(&example, id).Error
    return &example, err
}

func (r *exampleRepository) Create(example *storage.Example) (*storage.Example, error) {
    err := r.DB.DB.Create(example).Error
    return example, err
}

func (r *exampleRepository) Update(example *storage.Example) (*storage.Example, error) {
    err := r.DB.DB.Save(example).Error
    return example, err
}

func (r *exampleRepository) Delete(id uint) error {
    return r.DB.DB.Delete(&storage.Example{}, id).Error
}
```

### 4. Service (app/module/{name}/service/{name}_service.go)

```go
package service

import (
    "savory-ai-server/app/module/example/payload"
    "savory-ai-server/app/module/example/repository"
    "savory-ai-server/app/storage"
)

type exampleService struct {
    repo repository.ExampleRepository
}

type ExampleService interface {
    GetAll() ([]payload.ExampleResp, error)
    GetByID(id uint) (*payload.ExampleResp, error)
    Create(req *payload.ExampleCreateReq) (*payload.ExampleResp, error)
}

func NewExampleService(repo repository.ExampleRepository) ExampleService {
    return &exampleService{repo: repo}
}

func (s *exampleService) GetAll() ([]payload.ExampleResp, error) {
    examples, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    var resp []payload.ExampleResp
    for _, e := range examples {
        resp = append(resp, payload.ExampleResp{
            ID:   e.ID,
            Name: e.Name,
        })
    }
    return resp, nil
}
```

### 5. Controller (app/module/{name}/controller/{name}_controller.go)

```go
package controller

import (
    "github.com/gofiber/fiber/v2"
    "savory-ai-server/app/module/example/payload"
    "savory-ai-server/app/module/example/service"
    "savory-ai-server/utils/response"
)

type exampleController struct {
    service service.ExampleService
}

type ExampleController interface {
    GetAll(c *fiber.Ctx) error
    GetByID(c *fiber.Ctx) error
    Create(c *fiber.Ctx) error
}

func NewExampleController(service service.ExampleService) ExampleController {
    return &exampleController{service: service}
}

func (ctrl *exampleController) GetAll(ctx *fiber.Ctx) error {
    data, err := ctrl.service.GetAll()
    if err != nil {
        return err
    }

    return response.Resp(ctx, response.Response{
        Data:     data,
        Messages: response.Messages{"success"},
        Code:     fiber.StatusOK,
    })
}

func (ctrl *exampleController) Create(ctx *fiber.Ctx) error {
    req := new(payload.ExampleCreateReq)
    if err := ctx.BodyParser(req); err != nil {
        return err
    }

    // Валидация
    if err := response.ValidateStruct(req); err != nil {
        return response.Resp(ctx, response.Response{
            Messages: response.Messages{err.Error()},
            Code:     fiber.StatusBadRequest,
        })
    }

    data, err := ctrl.service.Create(req)
    if err != nil {
        return err
    }

    return response.Resp(ctx, response.Response{
        Data:     data,
        Messages: response.Messages{"created"},
        Code:     fiber.StatusCreated,
    })
}
```

### 6. Payload (app/module/{name}/payload/)

**request.go:**
```go
package payload

type ExampleCreateReq struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
}

type ExampleUpdateReq struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

**response.go:**
```go
package payload

import "time"

type ExampleResp struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 7. Controller aggregator (app/module/{name}/controller/controller.go)

```go
package controller

type Controller struct {
    Example ExampleController
}

func NewControllers(example ExampleController) *Controller {
    return &Controller{Example: example}
}
```

### 8. Module (app/module/{name}/{name}_module.go)

```go
package example

import (
    "github.com/gofiber/fiber/v2"
    "go.uber.org/fx"
    "savory-ai-server/app/module/example/controller"
    "savory-ai-server/app/module/example/repository"
    "savory-ai-server/app/module/example/service"
)

type ExampleRouter struct {
    App        fiber.Router
    Controller *controller.Controller
}

func NewExampleRouter(fiber *fiber.App, controller *controller.Controller) *ExampleRouter {
    return &ExampleRouter{
        App:        fiber,
        Controller: controller,
    }
}

var ExampleModule = fx.Options(
    fx.Provide(repository.NewExampleRepository),
    fx.Provide(service.NewExampleService),
    fx.Provide(controller.NewExampleController),
    fx.Provide(controller.NewControllers),
    fx.Provide(NewExampleRouter),
)

func (r *ExampleRouter) RegisterExampleRoutes(auth fiber.Handler) {
    ctrl := r.Controller.Example
    r.App.Route("/examples", func(router fiber.Router) {
        router.Get("/", ctrl.GetAll)
        router.Get("/:id", ctrl.GetByID)
        router.Post("/", auth, ctrl.Create)
    })
}
```

### 9. Зарегистрировать модуль в cmd/main.go

```go
import "savory-ai-server/app/module/example"

// В fx.New() добавить:
example.ExampleModule,
```

### 10. Зарегистрировать роуты в app/router/api.go

```go
// В структуру Router добавить:
ExampleRouter *example.ExampleRouter

// В Register() добавить:
r.ExampleRouter.RegisterExampleRoutes(authMiddleware)
```

## Соглашения по коду

### Именование

| Что | Формат | Пример |
|-----|--------|--------|
| Пакеты | snake_case | `menu_category` |
| Файлы | snake_case | `user_controller.go` |
| Структуры | PascalCase | `UserService` |
| Интерфейсы | PascalCase | `UserRepository` |
| Методы | PascalCase | `FindByID` |
| Приватные поля | camelCase | `userRepo` |
| JSON поля | snake_case | `created_at` |
| Константы | PascalCase или UPPER_CASE | `DefaultTimeout` |

### Валидация запросов

Используй теги `validate` из go-playground/validator:

```go
type CreateReq struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Phone    string `json:"phone" validate:"omitempty,e164"`
    Age      int    `json:"age" validate:"omitempty,min=0,max=150"`
    Password string `json:"password" validate:"required,min=6"`
}
```

### Формат ответов API

```go
// Успешный ответ
response.Resp(ctx, response.Response{
    Data:     data,
    Messages: response.Messages{"success"},
    Code:     fiber.StatusOK,
})

// Ошибка валидации
response.Resp(ctx, response.Response{
    Messages: response.Messages{err.Error()},
    Code:     fiber.StatusBadRequest,
})

// Созданный ресурс
response.Resp(ctx, response.Response{
    Data:     created,
    Messages: response.Messages{"created successfully"},
    Code:     fiber.StatusCreated,
})
```

### JWT авторизация

Получение данных текущего пользователя в контроллере:

```go
func (ctrl *exampleController) Create(ctx *fiber.Ctx) error {
    currentUser := ctx.Locals("user").(jwt.JWTData)
    companyID := currentUser.CompanyID
    userID := currentUser.ID
    // ...
}
```

## GORM паттерны

### Preload связей

```go
// Загрузить все связи
db.Preload(clause.Associations).First(&user, id)

// Загрузить конкретную связь
db.Preload("Company").First(&user, id)

// Вложенный preload
db.Preload("Company.Address").First(&user, id)
```

### Фильтрация

```go
// Where
db.Where("status = ?", "active").Find(&items)

// Multiple conditions
db.Where("status = ? AND company_id = ?", "active", companyID).Find(&items)

// In clause
db.Where("id IN ?", ids).Find(&items)
```

## Модули проекта

| Модуль | Описание | Основные endpoints |
|--------|----------|-------------------|
| auth | Аутентификация | POST /auth/login, POST /auth/register |
| user | Пользователи | GET/POST/PATCH /user |
| organization | Организации | /organizations |
| restaurant | Рестораны | /restaurants |
| menu_category | Категории меню | /menu-categories |
| dish | Блюда | /dishes |
| table | Столы | /tables |
| question | Вопросы (мультиязычные) | /questions |
| qr_code | QR-коды | /qr-codes |
| file_upload | Загрузка файлов | /upload |
| chat | Чат | /chat |

## Конфигурация

Файлы: `config/config.toml` (dev), `config/config.docker.toml` (docker)

Основные секции:
- `[app]` - порт, хост, таймауты
- `[db.postgres]` - DSN для PostgreSQL
- `[middleware.jwt]` - секрет и время жизни токена
- `[logger]` - настройки логирования

## Частые задачи

### Добавить новое поле в модель

1. Обновить структуру в `app/storage/{model}.go`
2. Запустить с флагом `-migrate`
3. Обновить payload request/response
4. Обновить service для маппинга полей

### Добавить новый endpoint

1. Добавить метод в интерфейс Controller
2. Реализовать метод в controller
3. Добавить роут в `{module}_module.go` → `RegisterRoutes()`

### Добавить middleware к роуту

```go
router.Get("/protected", auth, ctrl.ProtectedEndpoint)
router.Get("/public", ctrl.PublicEndpoint)
```

- Не запускай сервер без миграций в продакшене!
- Не запускай go build после того как внес изменения в зависимости или поменяешь код
- Не запускай go run после того как внес изменения в зависимости или поменяешь код
