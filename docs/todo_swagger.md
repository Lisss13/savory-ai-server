# План проверки Swagger спецификации

## Общая информация

**Файл:** `docs/swagger.yaml`
**Версия OpenAPI:** 3.0.3
**Дата создания плана:** 2025-11-28

---

## 1. Проверка соответствия путей (Paths)

### Auth Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `POST /auth/login` | `/auth/login` | ✅ OK | |
| `POST /auth/register` | `/auth/register` | ✅ OK | |
| `POST /auth/change-password` | `/auth/change-password` | ✅ OK | |
| `POST /auth/request-password-reset` | `/auth/request-password-reset` | ✅ OK | |
| `POST /auth/verify-password-reset` | `/auth/verify-password-reset` | ✅ OK | |
| `GET /auth/chek` (проверка токена) | `/auth/chek` | ⚠️ Проверить | Опечатка в "chek"? |

### User Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /user/:id` | `/user/{id}` | ✅ OK | |
| `PATCH /user/:id` | `/user/{id}` | ✅ OK | |
| `POST /user` | `/user` | ✅ OK | |

### Organization Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /organization` | `/organization` | ✅ OK | |
| `GET /organization/:id` | `/organization/{id}` | ✅ OK | |
| `PATCH /organization/:id` | `/organization/{id}` | ✅ OK | |
| `POST /organization/:id/users` | `/organization/{id}/users` | ✅ OK | |
| `DELETE /organization/:id/users` | `/organization/{id}/users` | ✅ OK | |
| `GET /organization/:id/languages` | `/organization/{id}/languages` | ✅ OK | |
| `POST /organization/:id/languages` | `/organization/{id}/languages` | ✅ OK | |
| `DELETE /organization/:id/languages` | `/organization/{id}/languages` | ✅ OK | |

### Languages Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /languages` | `/languages` | ✅ OK | |
| `GET /languages/:id` | `/languages/{id}` | ✅ OK | |
| `POST /languages` | `/languages` | ✅ OK | |
| `PATCH /languages/:id` | `/languages/{id}` | ✅ OK | |
| `DELETE /languages/:id` | `/languages/{id}` | ✅ OK | |

### Restaurant Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /restaurants` | `/restaurants` | ❓ Проверить | |
| `GET /restaurants/:id` | `/restaurants/{id}` | ❓ Проверить | |
| `GET /restaurants/organization/:organization_id` | `/restaurants/organization/{organization_id}` | ❓ Проверить | |
| `POST /restaurants` | `/restaurants` | ❓ Проверить | |
| `PUT /restaurants/:id` | `/restaurants/{id}` | ❓ Проверить | |
| `DELETE /restaurants/:id` | `/restaurants/{id}` | ❓ Проверить | |

### Table Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /tables` | `/tables` | ❓ Проверить | |
| `GET /tables/:id` | `/tables/{id}` | ❓ Проверить | |
| `GET /tables/restaurant/:restaurant_id` | `/tables/restaurant/{restaurant_id}` | ❓ Проверить | |
| `POST /tables` | `/tables` | ❓ Проверить | |
| `PUT /tables/:id` | `/tables/{id}` | ❓ Проверить | |
| `DELETE /tables/:id` | `/tables/{id}` | ❓ Проверить | |

### Menu Category Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /categories/restaurant/:restaurant_id` | `/categories/restaurant/{restaurant_id}` | ✅ OK | Публичный |
| `GET /categories/:id` | `/categories/{id}` | ✅ OK | Публичный |
| `POST /categories` | `/categories` | ❓ Проверить | Требует auth |
| `DELETE /categories/:id` | `/categories/{id}` | ❓ Проверить | Требует auth |

### Dish Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /dishes/restaurant/:restaurant_id` | `/dishes/restaurant/{restaurant_id}` | ✅ OK | Публичный |
| `GET /dishes/category/:restaurant_id` | `/dishes/category/{restaurant_id}` | ✅ OK | Публичный |
| `GET /dishes/dish-of-day/:restaurant_id` | `/dishes/dish-of-day/{restaurant_id}` | ✅ OK | Публичный |
| `POST /dishes/dish-of-day/:id` | `/dishes/dish-of-day/{id}` | ❓ Проверить | Требует auth |
| `GET /dishes/:id` | `/dishes/{id}` | ✅ OK | Публичный |
| `POST /dishes` | `/dishes` | ❓ Проверить | Требует auth |
| `PUT /dishes/:id` | `/dishes/{id}` | ❓ Проверить | Требует auth |
| `DELETE /dishes/:id` | `/dishes/{id}` | ❓ Проверить | Требует auth |

### Question Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /questions` | `/questions` | ❓ Проверить | |
| `GET /questions/language/:code` | `/questions/language/{code}` | ❓ Проверить | |
| `POST /questions` | `/questions` | ❓ Проверить | |
| `PUT /questions/:id` | `/questions/{id}` | ❓ Проверить | |
| `DELETE /questions/:id` | `/questions/{id}` | ❓ Проверить | |

### QR Code Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /qrcodes/restaurant/:restaurant_id` | | ❓ Проверить | |
| `GET /qrcodes/restaurant/:restaurant_id/download` | | ❓ Проверить | |
| `GET /qrcodes/restaurant/:restaurant_id/table/:table_id` | | ❓ Проверить | |
| `GET /qrcodes/restaurant/:restaurant_id/table/:table_id/download` | | ❓ Проверить | |

### File Upload Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `POST /uploads/images` | `/uploads/images` | ❓ Проверить | |

### Chat Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /chat/restaurant/:restaurant_id` | | ❓ Проверить | Legacy |
| `POST /chat/table/session/start` | | ❓ Проверить | |
| `POST /chat/table/session/close/:session_id` | | ❓ Проверить | |
| `POST /chat/table/message/send` | | ❓ Проверить | |
| `GET /chat/table/session/:session_id/messages` | | ❓ Проверить | |
| `GET /chat/table/session/:table_id` | | ❓ Проверить | |
| `POST /chat/restaurant/session/start` | | ❓ Проверить | |
| `POST /chat/restaurant/session/close/:session_id` | | ❓ Проверить | |
| `POST /chat/restaurant/message/send` | | ❓ Проверить | |
| `GET /chat/restaurant/session/:session_id/messages` | | ❓ Проверить | |
| `GET /chat/restaurant/sessions/:restaurant_id` | | ❓ Проверить | |

### Reservation Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /reservations/available/:restaurant_id` | `/reservations/available/{restaurant_id}` | ❓ Проверить | Публичный |
| `GET /reservations/my` | `/reservations/my` | ❓ Проверить | Публичный (query: phone) |
| `POST /reservations` | `/reservations` | ❓ Проверить | Публичный |
| `POST /reservations/:id/cancel/public` | `/reservations/{id}/cancel/public` | ❓ Проверить | Публичный |
| `GET /reservations` | `/reservations` | ❓ Проверить | Требует auth |
| `GET /reservations/:id` | `/reservations/{id}` | ❓ Проверить | Требует auth |
| `GET /reservations/restaurant/:restaurant_id` | `/reservations/restaurant/{restaurant_id}` | ❓ Проверить | Требует auth |
| `PATCH /reservations/:id` | `/reservations/{id}` | ❓ Проверить | Требует auth |
| `POST /reservations/:id/cancel` | `/reservations/{id}/cancel` | ❓ Проверить | Требует auth |
| `DELETE /reservations/:id` | `/reservations/{id}` | ❓ Проверить | Требует auth |

### Subscription Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /subscriptions` | `/subscriptions` | ❓ Проверить | |
| `GET /subscriptions/:id` | `/subscriptions/{id}` | ❓ Проверить | |
| `GET /subscriptions/organization/:organizationId` | `/subscriptions/organization/{organizationId}` | ❓ Проверить | |
| `GET /subscriptions/organization/:organizationId/active` | `/subscriptions/organization/{organizationId}/active` | ❓ Проверить | |
| `POST /subscriptions` | `/subscriptions` | ❓ Проверить | |
| `PUT /subscriptions/:id` | `/subscriptions/{id}` | ❓ Проверить | |
| `POST /subscriptions/:id/extend` | `/subscriptions/{id}/extend` | ❓ Проверить | |
| `POST /subscriptions/:id/deactivate` | `/subscriptions/{id}/deactivate` | ❓ Проверить | |
| `DELETE /subscriptions/:id` | `/subscriptions/{id}` | ❓ Проверить | |

### Admin Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /admin/stats` | `/admin/stats` | ❓ Проверить | Требует role: admin |
| `GET /admin/users` | `/admin/users` | ❓ Проверить | |
| `GET /admin/users/:id` | `/admin/users/{id}` | ❓ Проверить | |
| `PATCH /admin/users/:id/status` | `/admin/users/{id}/status` | ❓ Проверить | |
| `PATCH /admin/users/:id/role` | `/admin/users/{id}/role` | ❓ Проверить | |
| `DELETE /admin/users/:id` | `/admin/users/{id}` | ❓ Проверить | |
| `GET /admin/organizations` | `/admin/organizations` | ❓ Проверить | |
| `GET /admin/organizations/:id` | `/admin/organizations/{id}` | ❓ Проверить | |
| `DELETE /admin/organizations/:id` | `/admin/organizations/{id}` | ❓ Проверить | |
| `GET /admin/dishes` | `/admin/dishes` | ❓ Проверить | |
| `DELETE /admin/dishes/:id` | `/admin/dishes/{id}` | ❓ Проверить | |
| `GET /admin/logs` | `/admin/logs` | ❓ Проверить | |
| `GET /admin/logs/me` | `/admin/logs/me` | ❓ Проверить | |

---

## 2. Проверка схем (Schemas)

### Общие схемы
- [ ] `Response` - базовая обёртка ответов
- [ ] `ErrorResponse` - ответ при ошибке

### Auth
- [ ] `LoginRequest` - соответствует payload/request.go
- [ ] `RegisterRequest` - соответствует payload/request.go
- [ ] `ChangePasswordRequest` - соответствует payload/request.go
- [ ] `RequestPasswordResetRequest` - соответствует payload/request.go
- [ ] `VerifyPasswordResetRequest` - соответствует payload/request.go
- [ ] `LoginResponse` - соответствует payload/response.go

### User
- [ ] `UserResponse` - соответствует payload/response.go
- [ ] `UserCreateRequest` - соответствует payload/request.go
- [ ] `UserUpdateRequest` - соответствует payload/request.go

### Organization
- [ ] `OrganizationResponse` - соответствует payload/response.go
- [ ] `UpdateOrganizationRequest` - соответствует payload/request.go
- [ ] `AddUserToOrgRequest` - соответствует payload/request.go
- [ ] `RemoveUserFromOrgRequest` - соответствует payload/request.go

### Language
- [ ] `LanguageResponse` - соответствует payload/response.go
- [ ] `CreateLanguageRequest` - соответствует payload/request.go
- [ ] `UpdateLanguageRequest` - соответствует payload/request.go

### Restaurant
- [ ] `RestaurantResponse` - соответствует payload/response.go
- [ ] `CreateRestaurantRequest` - соответствует payload/request.go
- [ ] `UpdateRestaurantRequest` - соответствует payload/request.go
- [ ] `WorkingHourResponse` - соответствует payload/response.go
- [ ] `WorkingHourRequest` - соответствует payload/request.go

### Table
- [ ] `TableResponse` - соответствует payload/response.go
- [ ] `CreateTableRequest` - соответствует payload/request.go
- [ ] `UpdateTableRequest` - соответствует payload/request.go

### Menu Category
- [ ] `MenuCategoryResponse` - проверить наличие `restaurantId`
- [ ] `CreateMenuCategoryRequest` - проверить наличие `restaurant_id`

### Dish
- [ ] `DishResponse` - проверить `restaurant` вместо `organization`
- [ ] `CreateDishRequest` - проверить `restaurant_id`, `menu_category_id`, `allergens`
- [ ] `UpdateDishRequest` - проверить `restaurant_id`, `menu_category_id`, `allergens`
- [ ] `IngredientResponse` - соответствует payload/response.go
- [ ] `IngredientRequest` - соответствует payload/request.go
- [ ] `AllergenResponse` - соответствует payload/response.go
- [ ] `AllergenRequest` - соответствует payload/request.go

### Question
- [ ] `QuestionResponse` - проверить наличие `chat_type`
- [ ] `CreateQuestionRequest` - проверить наличие `chatType`
- [ ] `UpdateQuestionRequest` - проверить наличие `chatType`

### Chat
- [ ] `StartTableSessionRequest` - соответствует payload/request.go
- [ ] `SendTableMessageRequest` - соответствует payload/request.go
- [ ] `StartRestaurantSessionRequest` - соответствует payload/request.go
- [ ] `SendRestaurantMessageRequest` - соответствует payload/request.go
- [ ] `ChatMessageResponse` - соответствует payload/response.go
- [ ] `ChatSessionResponse` - соответствует payload/response.go

### Reservation
- [ ] `ReservationResponse` - соответствует payload/response.go
- [ ] `CreateReservationRequest` - соответствует payload/request.go
- [ ] `UpdateReservationRequest` - соответствует payload/request.go
- [ ] `AvailableSlotResponse` - соответствует payload/response.go
- [ ] `AvailableSlotsResponse` - соответствует payload/response.go
- [ ] `CancelReservationByPhoneRequest` - соответствует payload/request.go

### Subscription
- [ ] `SubscriptionResponse` - соответствует payload/response.go
- [ ] `CreateSubscriptionRequest` - соответствует payload/request.go
- [ ] `UpdateSubscriptionRequest` - соответствует payload/request.go
- [ ] `ExtendSubscriptionRequest` - соответствует payload/request.go

### Admin
- [ ] `AdminStatsResponse` - соответствует payload/response.go
- [ ] `AdminUserResponse` - соответствует payload/response.go
- [ ] `AdminUsersListResponse` - соответствует payload/response.go
- [ ] `AdminOrganizationResponse` - соответствует payload/response.go
- [ ] `AdminOrganizationsListResponse` - соответствует payload/response.go
- [ ] `AdminDishResponse` - проверить `restaurantId` вместо `organizationId`
- [ ] `AdminDishesListResponse` - соответствует payload/response.go
- [ ] `AdminLogResponse` - соответствует payload/response.go
- [ ] `AdminLogsListResponse` - соответствует payload/response.go
- [ ] `UpdateUserStatusRequest` - соответствует payload/request.go
- [ ] `UpdateUserRoleRequest` - соответствует payload/request.go

---

## 3. Проверка аутентификации

### Публичные эндпоинты (без auth)
- [ ] `GET /ping` - Health check
- [ ] `POST /auth/login` - Вход
- [ ] `POST /auth/register` - Регистрация
- [ ] `POST /auth/request-password-reset` - Запрос сброса пароля
- [ ] `POST /auth/verify-password-reset` - Подтверждение сброса
- [ ] `GET /categories/restaurant/{restaurant_id}` - Категории ресторана
- [ ] `GET /categories/{id}` - Категория по ID
- [ ] `GET /dishes/restaurant/{restaurant_id}` - Блюда ресторана
- [ ] `GET /dishes/category/{restaurant_id}` - Блюда по категориям
- [ ] `GET /dishes/dish-of-day/{restaurant_id}` - Блюдо дня
- [ ] `GET /dishes/{id}` - Блюдо по ID
- [ ] `GET /reservations/available/{restaurant_id}` - Доступные слоты
- [ ] `GET /reservations/my` - Мои бронирования
- [ ] `POST /reservations` - Создать бронирование
- [ ] `POST /reservations/{id}/cancel/public` - Отменить бронирование

### Защищённые эндпоинты (требуют auth)
- [ ] Все остальные эндпоинты должны иметь `security: BearerAuth`

### Admin-only эндпоинты
- [ ] Все `/admin/*` эндпоинты должны проверять роль admin

---

## 4. Проверка форматов данных

### JSON naming conventions
- [ ] Request/Response поля используют `snake_case` где нужно
- [ ] Проверить консистентность между `camelCase` и `snake_case`

### Типы данных
- [ ] Даты в формате `date-time` (ISO 8601)
- [ ] Цены в формате `number` с `format: float`
- [ ] ID в формате `integer`

---

## 5. Известные проблемы (исправлены ранее)

### Исправлено в текущей сессии:
1. ✅ `MenuCategoryResponse` - добавлено `restaurantId`
2. ✅ `CreateMenuCategoryRequest` - добавлено `restaurant_id`
3. ✅ `AllergenResponse` и `AllergenRequest` - добавлены схемы
4. ✅ `DishResponse` - заменено `organization` на `restaurant`
5. ✅ `CreateDishRequest` - добавлены `restaurant_id`, `menu_category_id`, `allergens`
6. ✅ `UpdateDishRequest` - добавлены `restaurant_id`, `menu_category_id`, `allergens`
7. ✅ `AdminDishResponse` - заменены `organizationId/Name` на `restaurantId/Name`
8. ✅ Путь `/categories` изменён на `/categories/restaurant/{restaurant_id}`
9. ✅ Путь `/dishes` изменён на `/dishes/restaurant/{restaurant_id}`
10. ✅ Путь `/dishes/category` изменён на `/dishes/category/{restaurant_id}`
11. ✅ Путь `/dishes/dish-of-day` изменён на `/dishes/dish-of-day/{restaurant_id}`
12. ✅ Публичные GET эндпоинты для блюд и категорий - убрана security

---

## 6. Рекомендуемый порядок проверки

1. **Сравнить пути** - пройти по каждому модулю и проверить все эндпоинты
2. **Проверить схемы** - сравнить swagger schemas с Go структурами в payload/
3. **Проверить аутентификацию** - убедиться что публичные/защищённые эндпоинты правильно отмечены
4. **Проверить примеры** - добавить примеры запросов/ответов где отсутствуют
5. **Валидировать спецификацию** - использовать swagger-editor для проверки синтаксиса

---

## 7. Инструменты для проверки

- **Swagger Editor**: https://editor.swagger.io/
- **OpenAPI Validator**: `npx swagger-cli validate docs/swagger.yaml`
- **Postman**: Импортировать swagger.yaml и протестировать эндпоинты

---

## Прогресс

| Раздел | Статус |
|--------|--------|
| Auth | ⏳ Требует проверки |
| User | ⏳ Требует проверки |
| Organization | ⏳ Требует проверки |
| Language | ⏳ Требует проверки |
| Restaurant | ⏳ Требует проверки |
| Table | ⏳ Требует проверки |
| Menu Category | ✅ Проверено |
| Dish | ✅ Проверено |
| Question | ⏳ Требует проверки |
| QR Code | ⏳ Требует проверки |
| File Upload | ⏳ Требует проверки |
| Chat | ⏳ Требует проверки |
| Reservation | ⏳ Требует проверки |
| Subscription | ⏳ Требует проверки |
| Admin | ✅ Проверено (схемы) |
