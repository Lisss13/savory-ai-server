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
| `GET /restaurants` | `/restaurants` | ✅ OK | auth |
| `GET /restaurants/:id` | `/restaurants/{id}` | ✅ OK | auth |
| `GET /restaurants/organization/:organization_id` | `/restaurants/organization/{organization_id}` | ✅ OK | auth |
| `POST /restaurants` | `/restaurants` | ✅ OK | auth |
| `PUT /restaurants/:id` | `/restaurants/{id}` | ✅ OK | auth |
| `DELETE /restaurants/:id` | `/restaurants/{id}` | ✅ OK | auth |

### Table Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /tables` | `/tables` | ✅ OK | auth |
| `GET /tables/:id` | `/tables/{id}` | ✅ OK | auth |
| `GET /tables/restaurant/:restaurant_id` | `/tables/restaurant/{restaurant_id}` | ✅ OK | auth |
| `POST /tables` | `/tables` | ✅ OK | auth |
| `PUT /tables/:id` | `/tables/{id}` | ✅ OK | auth |
| `DELETE /tables/:id` | `/tables/{id}` | ✅ OK | auth |

### Menu Category Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /categories/restaurant/:restaurant_id` | `/categories/restaurant/{restaurant_id}` | ✅ OK | Публичный |
| `GET /categories/:id` | `/categories/{id}` | ✅ OK | Публичный |
| `POST /categories` | `/categories` | ✅ OK | auth |
| `DELETE /categories/:id` | `/categories/{id}` | ✅ OK | auth |

### Dish Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /dishes/restaurant/:restaurant_id` | `/dishes/restaurant/{restaurant_id}` | ✅ OK | Публичный |
| `GET /dishes/category/:restaurant_id` | `/dishes/category/{restaurant_id}` | ✅ OK | Публичный |
| `GET /dishes/dish-of-day/:restaurant_id` | `/dishes/dish-of-day/{restaurant_id}` | ✅ OK | Публичный |
| `POST /dishes/dish-of-day/:id` | `/dishes/dish-of-day/{id}` | ✅ OK | auth |
| `GET /dishes/:id` | `/dishes/{id}` | ✅ OK | Публичный |
| `POST /dishes` | `/dishes` | ✅ OK | auth |
| `PUT /dishes/:id` | `/dishes/{id}` | ✅ OK | auth |
| `DELETE /dishes/:id` | `/dishes/{id}` | ✅ OK | auth |

### Question Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /questions` | `/questions` | ✅ OK | auth |
| `GET /questions/language/:code` | `/questions/language/{code}` | ✅ OK | auth |
| `POST /questions` | `/questions` | ✅ OK | auth |
| `PUT /questions/:id` | `/questions/{id}` | ✅ OK | auth |
| `DELETE /questions/:id` | `/questions/{id}` | ✅ OK | auth |

### QR Code Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /qrcodes/restaurant/:restaurant_id` | `/qrcodes/restaurant/{restaurant_id}` | ✅ OK | Публичный |
| `GET /qrcodes/restaurant/:restaurant_id/download` | `/qrcodes/restaurant/{restaurant_id}/download` | ✅ OK | Публичный |
| `GET /qrcodes/restaurant/:restaurant_id/table/:table_id` | `/qrcodes/restaurant/{restaurant_id}/table/{table_id}` | ✅ OK | Публичный |
| `GET /qrcodes/restaurant/:restaurant_id/table/:table_id/download` | `/qrcodes/restaurant/{restaurant_id}/table/{table_id}/download` | ✅ OK | Публичный |

### File Upload Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `POST /uploads/images` | `/uploads/images` | ✅ OK | Публичный |

### Chat Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /chat/restaurant/:restaurant_id` | `/chat/restaurant/{restaurant_id}` | ✅ OK | Legacy, deprecated |
| `POST /chat/table/session/start` | `/chat/table/session/start` | ✅ OK | Публичный |
| `POST /chat/table/session/close/:session_id` | `/chat/table/session/close/{session_id}` | ✅ OK | Публичный |
| `POST /chat/table/message/send` | `/chat/table/message/send` | ✅ OK | Публичный |
| `GET /chat/table/session/:session_id/messages` | `/chat/table/session/{session_id}/messages` | ✅ OK | Публичный |
| `GET /chat/table/session/:table_id` | `/chat/table/session/{table_id}` | ✅ OK | Публичный |
| `POST /chat/restaurant/session/start` | `/chat/restaurant/session/start` | ✅ OK | Публичный |
| `POST /chat/restaurant/session/close/:session_id` | `/chat/restaurant/session/close/{session_id}` | ✅ OK | Публичный |
| `POST /chat/restaurant/message/send` | `/chat/restaurant/message/send` | ✅ OK | Публичный |
| `GET /chat/restaurant/session/:session_id/messages` | `/chat/restaurant/session/{session_id}/messages` | ✅ OK | Публичный |
| `GET /chat/restaurant/sessions/:restaurant_id` | `/chat/restaurant/sessions/{restaurant_id}` | ✅ OK | Публичный |

### Reservation Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /reservations/available/:restaurant_id` | `/reservations/available/{restaurant_id}` | ✅ OK | Публичный |
| `GET /reservations/my` | `/reservations/my` | ✅ OK | Публичный (query: phone) |
| `POST /reservations` | `/reservations` | ✅ OK | Публичный |
| `POST /reservations/:id/cancel/public` | `/reservations/{id}/cancel/public` | ✅ OK | Публичный |
| `GET /reservations` | `/reservations` | ✅ OK | auth |
| `GET /reservations/:id` | `/reservations/{id}` | ✅ OK | auth |
| `GET /reservations/restaurant/:restaurant_id` | `/reservations/restaurant/{restaurant_id}` | ✅ OK | auth |
| `PATCH /reservations/:id` | `/reservations/{id}` | ✅ OK | auth |
| `POST /reservations/:id/cancel` | `/reservations/{id}/cancel` | ✅ OK | auth |
| `DELETE /reservations/:id` | `/reservations/{id}` | ✅ OK | auth |

### Subscription Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /subscriptions` | `/subscriptions` | ✅ OK | auth |
| `GET /subscriptions/:id` | `/subscriptions/{id}` | ✅ OK | auth |
| `GET /subscriptions/organization/:organizationId` | `/subscriptions/organization/{organizationId}` | ✅ OK | auth |
| `GET /subscriptions/organization/:organizationId/active` | `/subscriptions/organization/{organizationId}/active` | ✅ OK | auth |
| `POST /subscriptions` | `/subscriptions` | ✅ OK | auth |
| `PUT /subscriptions/:id` | `/subscriptions/{id}` | ✅ OK | auth |
| `POST /subscriptions/:id/extend` | `/subscriptions/{id}/extend` | ✅ OK | auth |
| `POST /subscriptions/:id/deactivate` | `/subscriptions/{id}/deactivate` | ✅ OK | auth |
| `DELETE /subscriptions/:id` | `/subscriptions/{id}` | ✅ OK | auth |

### Admin Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `GET /admin/stats` | `/admin/stats` | ✅ OK | admin |
| `GET /admin/users` | `/admin/users` | ✅ OK | admin |
| `GET /admin/users/:id` | `/admin/users/{id}` | ✅ OK | admin |
| `PATCH /admin/users/:id/status` | `/admin/users/{id}/status` | ✅ OK | admin |
| `PATCH /admin/users/:id/role` | `/admin/users/{id}/role` | ✅ OK | admin |
| `DELETE /admin/users/:id` | `/admin/users/{id}` | ✅ OK | admin |
| `GET /admin/organizations` | `/admin/organizations` | ✅ OK | admin |
| `GET /admin/organizations/:id` | `/admin/organizations/{id}` | ✅ OK | admin |
| `DELETE /admin/organizations/:id` | `/admin/organizations/{id}` | ✅ OK | admin |
| `GET /admin/dishes` | `/admin/dishes` | ✅ OK | admin |
| `DELETE /admin/dishes/:id` | `/admin/dishes/{id}` | ✅ OK | admin |
| `GET /admin/logs` | `/admin/logs` | ✅ OK | admin |
| `GET /admin/logs/me` | `/admin/logs/me` | ✅ OK | admin |

### Support Module
| Код | Swagger | Статус | Примечание |
|-----|---------|--------|------------|
| `POST /support` | `/support` | ✅ OK | auth |
| `GET /support/my` | `/support/my` | ✅ OK | auth |
| `GET /support/:id` | `/support/{id}` | ✅ OK | auth |
| `GET /admin/support` | `/admin/support` | ✅ OK | admin |
| `PATCH /admin/support/:id/status` | `/admin/support/{id}/status` | ✅ OK | admin |

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

### Support
- [ ] `SupportTicketResponse` - соответствует payload/response.go (id, user_id, user_name, user_email, title, description, email, phone, status, created_at, updated_at)
- [ ] `SupportTicketsListResponse` - соответствует payload/response.go (tickets[], total_count, page, page_size)
- [ ] `CreateSupportTicketRequest` - соответствует payload/request.go (title, description, email, phone)
- [ ] `UpdateSupportTicketStatusRequest` - соответствует payload/request.go (status: in_progress | completed)

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

### Support эндпоинты (требуют auth)
- [ ] `POST /support` - Создать заявку
- [ ] `GET /support/my` - Мои заявки
- [ ] `GET /support/{id}` - Заявка по ID

### Admin-only эндпоинты
- [ ] Все `/admin/*` эндпоинты должны проверять роль admin
- [ ] `GET /admin/support` - Все заявки в поддержку
- [ ] `PATCH /admin/support/{id}/status` - Обновить статус заявки

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
| Auth | ✅ Проверено |
| User | ✅ Проверено |
| Organization | ✅ Проверено |
| Language | ✅ Проверено |
| Restaurant | ✅ Проверено |
| Table | ✅ Проверено |
| Menu Category | ✅ Проверено |
| Dish | ✅ Проверено |
| Question | ✅ Проверено |
| QR Code | ✅ Проверено |
| File Upload | ✅ Проверено |
| Chat | ✅ Проверено |
| Reservation | ✅ Проверено |
| Subscription | ✅ Проверено |
| Admin | ✅ Проверено |
| Support | ✅ Проверено |
