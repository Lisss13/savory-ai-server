# Рекомендации по улучшению Savory AI Server

> Анализ выполнен: 25 ноября 2025

## Содержание

1. [Критичные проблемы](#критичные-проблемы)
2. [Безопасность](#безопасность)
3. [Обработка ошибок](#обработка-ошибок)
4. [Валидация данных](#валидация-данных)
5. [Производительность](#производительность)
6. [Логирование](#логирование)
7. [Тестирование](#тестирование)
8. [Конфигурация](#конфигурация)
9. [Код контроллеров и сервисов](#код-контроллеров-и-сервисов)
10. [Приоритизированный план действий](#приоритизированный-план-действий)

---

## Критичные проблемы

### 1. Небезопасная генерация кодов сброса пароля

**Файл:** `app/module/auth/service/auth_service.go:199`

```go
// Текущий код (НЕБЕЗОПАСНО)
rand.Seed(time.Now().UnixNano())
code := strconv.Itoa(100000 + rand.Intn(900000))
```

**Проблемы:**
- Используется `math/rand` вместо `crypto/rand`
- `rand.Seed()` вызывается при каждом запросе
- Всего 900k комбинаций - легко перебрать

**Решение:**
```go
import "crypto/rand"

func generateSecureCode() string {
    b := make([]byte, 3)
    rand.Read(b)
    code := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
    return fmt.Sprintf("%06d", code%1000000)
}
```

### 2. Отсутствие проверки token.Valid в JWT

**Файл:** `utils/jwt/jwt.go:46-71`

```go
// Текущий код возвращает данные даже для невалидного токена
if claims, ok := token.Claims.(jwt.MapClaims); ok {
    // извлечение данных без проверки token.Valid
}
return token.Valid, data
```

**Решение:**
```go
func (j *JWT) ParseToken(tokenString string) (bool, *JWTData) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(j.Secret), nil
    })

    if err != nil || !token.Valid {
        return false, nil
    }

    // Теперь безопасно извлекать claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return false, nil
    }
    // ...
}
```

### 3. Отсутствие транзакций в multi-step операциях

**Файл:** `app/module/dish/repository/dish_repository.go:79-128`

```go
// Текущий код - без транзакции
func (r *dishRepository) Update(dish *storage.Dish) (*storage.Dish, error) {
    r.DB.DB.Model(&dish).Updates(...)     // Шаг 1
    r.DB.DB.Delete(&storage.Ingredient{}) // Шаг 2 - если падает, данные повреждены
    r.DB.DB.Create(&dish.Ingredients)     // Шаг 3
}
```

**Решение:**
```go
func (r *dishRepository) Update(dish *storage.Dish) (*storage.Dish, error) {
    err := r.DB.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&dish).Updates(...).Error; err != nil {
            return err
        }
        if err := tx.Where("dish_id = ?", dish.ID).Delete(&storage.Ingredient{}).Error; err != nil {
            return err
        }
        if len(dish.Ingredients) > 0 {
            if err := tx.Create(&dish.Ingredients).Error; err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return r.FindByID(dish.ID)
}
```

---

## Безопасность

### 4. Отсутствие CORS middleware

**Файл:** `app/middleware/register.go`

**Проблема:** API открыта для всех источников.

**Решение:**
```go
import "github.com/gofiber/fiber/v2/middleware/cors"

func (m *Middleware) Register() {
    m.App.Use(cors.New(cors.Config{
        AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"), // "https://app.savory.ai"
        AllowMethods:     "GET,POST,PATCH,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Content-Type,Authorization",
        AllowCredentials: true,
    }))
    // остальные middleware
}
```

### 5. Секреты в config.toml

**Файл:** `config/config.toml`

```toml
# НЕБЕЗОПАСНО - секреты в git!
[middleware.jwt]
secret = "/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="

[db.postgres]
dsn = "postgresql://savory_admin:savory_password@localhost:5432/savory_db"
```

**Решение:**
1. Добавить `config/config.toml` в `.gitignore`
2. Создать `config/config.example.toml` с placeholder'ами
3. Использовать переменные окружения:

```go
// utils/config/config.go
func NewConfig() *Config {
    config := parseConfigFile()

    // Override с env переменными
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        config.DB.Postgres.DSN = dbURL
    }
    if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
        config.Middleware.Jwt.Secret = jwtSecret
    }

    return config
}
```

### 6. Rate limiting для публичных endpoints

**Файл:** `app/module/reservation/reservation_module.go:61-62`

**Проблема:** Публичные endpoints не защищены от brute-force атак.

**Решение:**
```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

func (r *ReservationRouter) RegisterReservationRoutes(auth fiber.Handler) {
    // Rate limiter для публичных endpoints
    publicLimiter := limiter.New(limiter.Config{
        Max:        10,              // 10 запросов
        Expiration: 1 * time.Minute, // в минуту
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
    })

    r.App.Route("/reservations", func(router fiber.Router) {
        // Публичные с rate limiting
        router.Get("/available/:restaurant_id", publicLimiter, ctrl.GetAvailableSlots)
        router.Post("/", publicLimiter, ctrl.Create)

        // Защищённые
        router.Get("/", auth, ctrl.GetAll)
    })
}
```

---

## Обработка ошибок

### 7. Игнорирование ошибок в GetStats

**Файл:** `app/module/admin/service/admin_service.go:52-59`

```go
// Текущий код - все ошибки игнорируются
totalUsers, _ := s.adminRepo.CountUsers()
activeUsers, _ := s.adminRepo.CountActiveUsers()
totalOrgs, _ := s.adminRepo.CountOrganizations()
```

**Решение:**
```go
func (s *adminService) GetStats() (*payload.AdminStatsResp, error) {
    var errors []string

    totalUsers, err := s.adminRepo.CountUsers()
    if err != nil {
        log.Error().Err(err).Msg("Failed to count users")
        errors = append(errors, "users count unavailable")
    }

    // ... остальные запросы

    return &payload.AdminStatsResp{
        TotalUsers: totalUsers,
        Errors:     errors, // Добавить поле для ошибок
    }, nil
}
```

### 8. Создать helper для обработки ошибок в контроллерах

**Проблема:** Одинаковый код повторяется 50+ раз.

**Решение:** Создать `utils/response/helpers.go`:

```go
package response

import "github.com/gofiber/fiber/v2"

// HandleResult обрабатывает результат сервиса и возвращает ответ
func HandleResult(c *fiber.Ctx, data any, err error, successMsg string) error {
    if err != nil {
        return Resp(c, Response{
            Messages: Messages{err.Error()},
            Code:     fiber.StatusBadRequest,
        })
    }
    return Resp(c, Response{
        Data:     data,
        Messages: Messages{successMsg},
        Code:     fiber.StatusOK,
    })
}

// HandleCreated для успешного создания ресурса
func HandleCreated(c *fiber.Ctx, data any, err error) error {
    if err != nil {
        return Resp(c, Response{
            Messages: Messages{err.Error()},
            Code:     fiber.StatusBadRequest,
        })
    }
    return Resp(c, Response{
        Data:     data,
        Messages: Messages{"created successfully"},
        Code:     fiber.StatusCreated,
    })
}
```

**Использование в контроллерах:**
```go
func (ctrl *dishController) Create(ctx *fiber.Ctx) error {
    req := new(payload.CreateDishReq)
    if err := ctx.BodyParser(req); err != nil {
        return response.HandleResult(ctx, nil, err, "")
    }

    data, err := ctrl.service.Create(req)
    return response.HandleCreated(ctx, data, err)
}
```

---

## Валидация данных

### 9. Валидация прошедших дат в бронировании

**Файл:** `app/module/reservation/service/reservation_service.go:290-301`

```go
// Текущий код - не проверяет прошедшие даты
date, err := time.Parse("2006-01-02", req.ReservationDate)
if err != nil {
    return nil, errors.New("invalid date format")
}
// Можно создать бронь на вчера!
```

**Решение:**
```go
func (s *reservationService) Create(req *payload.CreateReservationReq) (*payload.ReservationResp, error) {
    date, err := time.Parse("2006-01-02", req.ReservationDate)
    if err != nil {
        return nil, errors.New("invalid date format, use YYYY-MM-DD")
    }

    // Проверка на прошедшую дату
    today := time.Now().Truncate(24 * time.Hour)
    if date.Before(today) {
        return nil, errors.New("cannot create reservation for past dates")
    }

    // Проверка на слишком далёкую дату (например, не более 90 дней)
    maxDate := today.AddDate(0, 3, 0)
    if date.After(maxDate) {
        return nil, errors.New("cannot create reservation more than 90 days in advance")
    }

    // ...
}
```

### 10. Добавить ограничения длины полей

**Файл:** `app/storage/user.go`

```go
// Добавить constraints в GORM и validate теги
type User struct {
    gorm.Model
    Name     string `gorm:"column:name;type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
    Company  string `gorm:"column:company;type:varchar(200);not null" json:"company" validate:"required,min=2,max=200"`
    Email    string `gorm:"column:email;type:varchar(255);unique;not null" json:"email" validate:"required,email,max=255"`
    Phone    string `gorm:"column:phone;type:varchar(20);not null" json:"phone" validate:"required,max=20"`
    Password string `gorm:"column:password;not null" json:"-"`
}
```

### 11. Валидация в Repository слое

**Все репозитории:**

```go
func (r *reservationRepository) FindByID(id uint) (*storage.Reservation, error) {
    if id == 0 {
        return nil, errors.New("invalid reservation id")
    }

    var reservation storage.Reservation
    err := r.DB.DB.Preload("Restaurant").Preload("Table").First(&reservation, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("reservation not found")
        }
        return nil, err
    }
    return &reservation, nil
}
```

---

## Производительность

### 12. Оптимизация N+1 запросов

**Файл:** `app/module/dish/repository/dish_repository.go:32`

```go
// Текущий код - 4 отдельных preload
err := r.DB.DB.Preload("Organization").Preload("MenuCategory").Preload("Ingredients").Preload("Allergens").Find(&dishes).Error
```

**Решение с использованием Joins:**
```go
// Для списка блюд - выбирать только нужные поля
func (r *dishRepository) FindAll() ([]storage.Dish, error) {
    var dishes []storage.Dish
    err := r.DB.DB.
        Select("dishes.*, menu_categories.name as category_name").
        Joins("LEFT JOIN menu_categories ON dishes.menu_category_id = menu_categories.id").
        Preload("Ingredients").
        Preload("Allergens").
        Find(&dishes).Error
    return dishes, err
}

// Для деталей блюда - полный preload
func (r *dishRepository) FindByID(id uint) (*storage.Dish, error) {
    var dish storage.Dish
    err := r.DB.DB.
        Preload(clause.Associations).
        First(&dish, id).Error
    return &dish, err
}
```

### 13. Преаллокация слайсов

**Файл:** `app/module/reservation/service/reservation_service.go:71-76`

```go
// Текущий код - append растёт динамически
var reservationResps []payload.ReservationResp
for _, res := range reservations {
    reservationResps = append(reservationResps, mapReservationToResponse(&res))
}
```

**Решение:**
```go
// Преаллокация для известного размера
reservationResps := make([]payload.ReservationResp, 0, len(reservations))
for _, res := range reservations {
    reservationResps = append(reservationResps, mapReservationToResponse(&res))
}

// Или ещё лучше
reservationResps := make([]payload.ReservationResp, len(reservations))
for i, res := range reservations {
    reservationResps[i] = mapReservationToResponse(&res)
}
```

### 14. Кэширование часто запрашиваемых данных

**Рекомендация:** Добавить кэширование для:
- Список категорий меню (редко меняется)
- Рабочие часы ресторана
- Список столиков ресторана

```go
// Пример с go-cache
import "github.com/patrickmn/go-cache"

type CachedRestaurantService struct {
    repo  repository.RestaurantRepository
    cache *cache.Cache
}

func (s *CachedRestaurantService) GetWorkingHours(restaurantID uint) ([]WorkingHour, error) {
    cacheKey := fmt.Sprintf("working_hours_%d", restaurantID)

    if cached, found := s.cache.Get(cacheKey); found {
        return cached.([]WorkingHour), nil
    }

    hours, err := s.repo.GetWorkingHours(restaurantID)
    if err != nil {
        return nil, err
    }

    s.cache.Set(cacheKey, hours, 5*time.Minute)
    return hours, nil
}
```

---

## Логирование

### 15. Добавить структурированное логирование в сервисы

**Файл:** `app/module/reservation/service/reservation_service.go`

```go
import "github.com/rs/zerolog/log"

func (s *reservationService) Create(req *payload.CreateReservationReq) (*payload.ReservationResp, error) {
    log.Info().
        Uint("restaurant_id", req.RestaurantID).
        Uint("table_id", req.TableID).
        Str("customer_phone", maskPhone(req.CustomerPhone)).
        Str("date", req.ReservationDate).
        Msg("Creating reservation")

    // ... логика ...

    if hasConflict {
        log.Warn().
            Uint("table_id", req.TableID).
            Str("time", req.StartTime).
            Msg("Slot conflict detected")
        return nil, errors.New("this time slot is already booked")
    }

    // После успешного создания
    log.Info().
        Uint("reservation_id", reservation.ID).
        Msg("Reservation created successfully")

    return result, nil
}

// Маскирование телефона для логов
func maskPhone(phone string) string {
    if len(phone) < 4 {
        return "***"
    }
    return phone[:3] + "****" + phone[len(phone)-2:]
}
```

### 16. Логирование аутентификации

**Файл:** `app/module/auth/service/auth_service.go`

```go
func (as *authService) Login(req payload.LoginRequest) (payload.LoginResponse, error) {
    user, err := as.userRepo.FindUserByEmail(req.Email)
    if err != nil {
        log.Warn().
            Str("email", req.Email).
            Str("ip", getClientIP()).
            Msg("Login attempt for non-existent user")
        return payload.LoginResponse{}, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword(...); err != nil {
        log.Warn().
            Str("email", req.Email).
            Uint("user_id", user.ID).
            Msg("Failed login attempt - wrong password")
        return payload.LoginResponse{}, errors.New("invalid credentials")
    }

    log.Info().
        Str("email", req.Email).
        Uint("user_id", user.ID).
        Msg("User logged in successfully")

    return response, nil
}
```

---

## Тестирование

### 17. Создать unit тесты для сервисов

**Создать:** `app/module/reservation/service/reservation_service_test.go`

```go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock репозитории
type MockReservationRepo struct {
    mock.Mock
}

func (m *MockReservationRepo) FindByID(id uint) (*storage.Reservation, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*storage.Reservation), args.Error(1)
}

// Тесты
func TestReservationService_Create_InvalidDate(t *testing.T) {
    mockRepo := new(MockReservationRepo)
    svc := NewReservationService(mockRepo, nil, nil)

    req := &payload.CreateReservationReq{
        ReservationDate: "2020-01-01", // Прошедшая дата
        RestaurantID:    1,
        TableID:         1,
        GuestCount:      2,
    }

    _, err := svc.Create(req)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "past dates")
}

func TestReservationService_Create_Success(t *testing.T) {
    mockRepo := new(MockReservationRepo)
    mockTableRepo := new(MockTableRepo)
    mockRestaurantRepo := new(MockRestaurantRepo)

    // Setup mocks
    mockRestaurantRepo.On("FindByID", uint(1)).Return(&storage.Restaurant{
        Model: gorm.Model{ID: 1},
        Name:  "Test Restaurant",
    }, nil)

    mockTableRepo.On("FindByID", uint(1)).Return(&storage.Table{
        Model:      gorm.Model{ID: 1},
        GuestCount: 4,
    }, nil)

    mockRepo.On("Create", mock.Anything).Return(&storage.Reservation{
        Model: gorm.Model{ID: 1},
    }, nil)

    svc := NewReservationService(mockRepo, mockTableRepo, mockRestaurantRepo)

    req := &payload.CreateReservationReq{
        ReservationDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
        StartTime:       "19:00",
        RestaurantID:    1,
        TableID:         1,
        GuestCount:      2,
        CustomerName:    "Test User",
        CustomerPhone:   "+79001234567",
    }

    result, err := svc.Create(req)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### 18. Добавить integration тесты

**Создать:** `tests/integration/reservation_test.go`

```go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func TestReservationAPI_Create(t *testing.T) {
    app := setupTestApp() // Настроить Fiber с тестовой БД

    payload := map[string]interface{}{
        "restaurant_id":    1,
        "table_id":         1,
        "customer_name":    "Test User",
        "customer_phone":   "+79001234567",
        "guest_count":      2,
        "reservation_date": "2025-12-01",
        "start_time":       "19:00",
    }

    body, _ := json.Marshal(payload)
    req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

---

## Конфигурация

### 19. Валидация URL в конфигурации

**Файл:** `utils/config/config.go`

```go
import "net/url"

func validateConfig(cfg *Config) {
    // Существующие проверки...

    // Валидация URL
    if cfg.App.ChatServiceUrl != "" {
        if _, err := url.ParseRequestURI(cfg.App.ChatServiceUrl); err != nil {
            log.Panic().Err(err).Str("url", cfg.App.ChatServiceUrl).Msg("Invalid Chat Service URL")
        }
    }

    // Валидация DSN
    if !strings.HasPrefix(cfg.DB.Postgres.DSN, "postgresql://") &&
       !strings.HasPrefix(cfg.DB.Postgres.DSN, "postgres://") {
        log.Panic().Msg("Invalid PostgreSQL DSN format")
    }
}
```

### 20. Graceful shutdown

**Файл:** `cmd/main.go`

```go
func main() {
    app := fx.New(
        // modules...
        fx.Invoke(func(lc fx.Lifecycle, fiber *fiber.App, cfg *config.Config) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    go func() {
                        if err := fiber.Listen(cfg.App.Port); err != nil {
                            log.Fatal().Err(err).Msg("Server failed to start")
                        }
                    }()
                    return nil
                },
                OnStop: func(ctx context.Context) error {
                    log.Info().Msg("Shutting down server...")
                    return fiber.ShutdownWithTimeout(30 * time.Second)
                },
            })
        }),
    )

    // Handle OS signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go app.Run()

    <-quit
    log.Info().Msg("Received shutdown signal")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := app.Stop(ctx); err != nil {
        log.Error().Err(err).Msg("Error during shutdown")
    }
}
```

---

## Код контроллеров и сервисов

### 21. Исправить непрофессиональные комментарии

**Файл:** `utils/response/response.go:33`

```go
// Было:
// nothiing to describe this fucking variable
var IsProduction bool

// Должно быть:
// IsProduction indicates whether the application runs in production mode.
// When true, error details are hidden from API responses.
var IsProduction bool
```

### 22. Унифицировать использование указателей в интерфейсах

**Файл:** `app/module/reservation/repository/reservation_repository.go`

```go
// Текущий код - несогласованно
type ReservationRepository interface {
    FindAll() ([]storage.Reservation, error)        // Слайс значений
    FindByID(id uint) (*storage.Reservation, error) // Указатель
}

// Рекомендация - всегда указатели для структур
type ReservationRepository interface {
    FindAll() ([]*storage.Reservation, error)
    FindByID(id uint) (*storage.Reservation, error)
    FindByPhone(phone string) ([]*storage.Reservation, error)
    Create(r *storage.Reservation) (*storage.Reservation, error)
}
```

### 23. Context timeout для внешних сервисов

**Файл:** `app/module/chat/service/chat_service.go`

```go
// Было: 30 секунд - слишком много
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

// Рекомендация: конфигурируемый timeout
type ChatService struct {
    aiTimeout time.Duration // из конфига
}

func (s *ChatService) ProcessMessage(...) {
    ctx, cancel := context.WithTimeout(context.Background(), s.aiTimeout)
    defer cancel()

    // Для пользовательского опыта лучше 10-15 секунд с retry
}
```

---

## Приоритизированный план действий

### Критичные (сделать немедленно)

| # | Проблема | Файл | Сложность |
|---|----------|------|-----------|
| 1 | Небезопасная генерация кодов | auth_service.go:199 | Низкая |
| 2 | Проверка token.Valid в JWT | jwt.go:46-71 | Низкая |
| 3 | Транзакции в Update операциях | dish_repository.go | Средняя |
| 4 | Секреты в config.toml | config.toml | Низкая |

### Высокий приоритет (в течение недели)

| # | Проблема | Файл | Сложность |
|---|----------|------|-----------|
| 5 | Unit тесты для сервисов | новые файлы | Высокая |
| 6 | Логирование в сервисах | все сервисы | Средняя |
| 7 | CORS middleware | middleware/register.go | Низкая |
| 8 | Rate limiting публичных endpoints | router | Средняя |

### Средний приоритет (в течение месяца)

| # | Проблема | Файл | Сложность |
|---|----------|------|-----------|
| 9 | Валидация дат в бронировании | reservation_service.go | Низкая |
| 10 | Helper для ошибок в контроллерах | response/helpers.go | Низкая |
| 11 | Ограничения длины полей | storage/*.go | Средняя |
| 12 | Обработка ошибок в GetStats | admin_service.go | Низкая |

### Низкий приоритет (backlog)

| # | Проблема | Файл | Сложность |
|---|----------|------|-----------|
| 13 | Оптимизация N+1 запросов | репозитории | Средняя |
| 14 | Кэширование | новые файлы | Высокая |
| 15 | Унификация указателей | интерфейсы | Средняя |
| 16 | Integration тесты | новые файлы | Высокая |

---

## Метрики успеха

После внедрения улучшений:

- [ ] Покрытие тестами service слоя > 70%
- [ ] Все секреты перенесены в env переменные
- [ ] Логирование всех критичных операций
- [ ] Нет игнорируемых ошибок (`_`) в критичном коде
- [ ] Rate limiting на всех публичных endpoints
- [ ] CORS настроен для production домена

---

## Полезные ссылки

- [OWASP Go Security](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [GORM Best Practices](https://gorm.io/docs/performance.html)
- [Fiber Security](https://docs.gofiber.io/api/middleware/helmet)
