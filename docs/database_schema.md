# Схема базы данных Savory AI

## ER-диаграмма

```mermaid
erDiagram
    %% =====================================================
    %% USERS & ORGANIZATIONS
    %% =====================================================

    User {
        uint id PK
        string name
        string company
        string email UK
        string phone
        string password
        UserRole role "user | admin"
        bool is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    Organization {
        uint id PK
        string name
        string phone
        uint admin_id FK
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    OrganizationUser {
        uint organization_id PK,FK
        uint user_id PK,FK
    }

    Language {
        uint id PK
        string code UK
        string name
        string description
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    OrganizationLanguage {
        uint organization_id PK,FK
        uint language_id PK,FK
    }

    Subscription {
        uint id PK
        uint organization_id FK
        int period "months"
        timestamp start_date
        timestamp end_date
        bool is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    PasswordResetCode {
        uint id PK
        uint user_id FK
        string code
        timestamp expires_at
        bool used
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    AdminLog {
        uint id PK
        uint admin_id FK
        AdminAction action
        string entity_type
        uint entity_id
        text details "JSON"
        string ip_address
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    %% =====================================================
    %% RESTAURANT
    %% =====================================================

    Restaurant {
        uint id PK
        uint organization_id FK
        string name
        string address
        string phone
        string website
        string description
        string image_url
        string menu
        int reservation_duration "minutes, default 90"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    WorkingHour {
        uint id PK
        uint restaurant_id FK
        int day_of_week "0-6, 0=Sunday"
        string open_time "HH:MM"
        string close_time "HH:MM"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    Table {
        uint id PK
        uint restaurant_id FK
        string name
        int guest_count
        string menu
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    %% =====================================================
    %% RESERVATIONS
    %% =====================================================

    Reservation {
        uint id PK
        uint restaurant_id FK
        uint table_id FK
        string customer_name
        string customer_phone
        string customer_email
        int guest_count
        date reservation_date
        string start_time "HH:MM"
        string end_time "HH:MM"
        ReservationStatus status "pending|confirmed|cancelled|completed|no_show"
        string notes
        uint chat_session_id FK "nullable"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    %% =====================================================
    %% MENU
    %% =====================================================

    MenuCategory {
        uint id PK
        uint restaurant_id FK
        string name
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    Dish {
        uint id PK
        uint restaurant_id FK
        uint menu_category_id FK
        string name
        float price
        string description
        string image
        bool is_dish_of_day
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    Ingredient {
        uint id PK
        uint dish_id FK
        string name
        float quantity
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    Allergen {
        uint id PK
        uint dish_id FK
        string name
        string description
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    %% =====================================================
    %% CHAT
    %% =====================================================

    Question {
        uint id PK
        uint organization_id FK
        string text
        uint language_id FK "nullable"
        ChatType chat_type "reservation | menu"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    TableChatSession {
        uint id PK
        uint table_id FK
        uint restaurant_id FK
        bool active
        timestamp last_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    TableChatMessage {
        uint id PK
        uint table_id FK
        uint restaurant_id FK
        uint chat_session_id FK
        string author_type "user | bot | restaurant"
        string content
        timestamp sent_at
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    RestaurantChatSession {
        uint id PK
        uint restaurant_id FK
        bool active
        timestamp last_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    RestaurantChatMessage {
        uint id PK
        uint restaurant_id FK
        uint chat_session_id FK
        string author_type "user | bot | restaurant"
        string content
        timestamp sent_at
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    %% =====================================================
    %% RELATIONSHIPS
    %% =====================================================

    %% User & Organization
    User ||--o{ Organization : "admin_id (администрирует)"
    Organization ||--o{ OrganizationUser : ""
    User ||--o{ OrganizationUser : ""

    %% Languages
    Organization ||--o{ OrganizationLanguage : ""
    Language ||--o{ OrganizationLanguage : ""

    %% Subscription
    Organization ||--o| Subscription : "has"

    %% Password Reset
    User ||--o{ PasswordResetCode : "has"

    %% Admin Logs
    User ||--o{ AdminLog : "admin_id"

    %% Restaurant
    Organization ||--o{ Restaurant : "owns"
    Restaurant ||--o{ WorkingHour : "has"
    Restaurant ||--o{ Table : "has"

    %% Reservations
    Restaurant ||--o{ Reservation : "has"
    Table ||--o{ Reservation : "has"

    %% Menu
    Restaurant ||--o{ MenuCategory : "has"
    Restaurant ||--o{ Dish : "has"
    MenuCategory ||--o{ Dish : "contains"
    Dish ||--o{ Ingredient : "has"
    Dish ||--o{ Allergen : "has"

    %% Questions
    Organization ||--o{ Question : "has"
    Language ||--o{ Question : "language_id"

    %% Table Chat
    Table ||--o{ TableChatSession : "has"
    Restaurant ||--o{ TableChatSession : "has"
    TableChatSession ||--o{ TableChatMessage : "has"
    Table ||--o{ TableChatMessage : ""
    Restaurant ||--o{ TableChatMessage : ""

    %% Restaurant Chat
    Restaurant ||--o{ RestaurantChatSession : "has"
    RestaurantChatSession ||--o{ RestaurantChatMessage : "has"
    Restaurant ||--o{ RestaurantChatMessage : ""
```

## Таблицы

### Пользователи и организации

| Таблица | Описание |
|---------|----------|
| `users` | Пользователи системы (роли: user, admin) |
| `organizations` | Организации (рестораннные сети) |
| `organization_users` | Many-to-Many связь пользователей и организаций |
| `languages` | Языки (en, ru и др.) |
| `organization_languages` | Many-to-Many связь организаций и языков |
| `subscriptions` | Подписки организаций |
| `password_reset_codes` | Коды сброса пароля |
| `admin_logs` | Логи действий администраторов |

### Рестораны

| Таблица | Описание |
|---------|----------|
| `restaurants` | Рестораны организации |
| `working_hours` | Рабочие часы ресторана по дням недели |
| `tables` | Столы в ресторане |

### Бронирования

| Таблица | Описание |
|---------|----------|
| `reservations` | Бронирования столов |

**Статусы бронирования:**
- `pending` — ожидает подтверждения
- `confirmed` — подтверждено
- `cancelled` — отменено
- `completed` — завершено
- `no_show` — гость не пришёл

### Меню

| Таблица | Описание |
|---------|----------|
| `menu_categories` | Категории меню |
| `dishes` | Блюда |
| `ingredients` | Ингредиенты блюд |
| `allergens` | Аллергены в блюдах |

### Чат

| Таблица | Описание |
|---------|----------|
| `questions` | Быстрые вопросы для чат-бота |
| `table_chat_sessions` | Сессии чата для стола (QR-код) |
| `table_chat_messages` | Сообщения чата стола |
| `restaurant_chat_sessions` | Сессии чата для бронирования |
| `restaurant_chat_messages` | Сообщения чата бронирования |

**Типы авторов сообщений:**
- `user` — пользователь/гость
- `bot` — AI-бот
- `restaurant` — сотрудник ресторана

**Типы чата (для вопросов):**
- `reservation` — чат бронирования
- `menu` — чат меню

## Индексы

Автоматически создаваемые GORM индексы:
- `id` (PK) на всех таблицах
- `users.email` (UNIQUE)
- `languages.code` (UNIQUE)
- `reservations.restaurant_id` (INDEX)
- `reservations.table_id` (INDEX)
- `reservations.reservation_date` (INDEX)

---

## UML Class Diagram (PlantUML)

```plantuml
@startuml Savory AI Database Schema

!theme plain
skinparam linetype ortho
skinparam classAttributeIconSize 0
skinparam classFontSize 12
skinparam packageFontSize 14

' =====================================================
' ENUMS
' =====================================================
package "Enums" <<Rectangle>> {
    enum UserRole {
        user
        admin
    }

    enum ReservationStatus {
        pending
        confirmed
        cancelled
        completed
        no_show
    }

    enum ChatType {
        reservation
        menu
    }

    enum AdminAction {
        create
        update
        delete
        block
        unblock
        activate
        deactivate
        view
    }

    enum AuthorType {
        user
        bot
        restaurant
    }
}

' =====================================================
' USERS & ORGANIZATIONS
' =====================================================
package "Users & Organizations" <<Frame>> {
    class User {
        +id : uint <<PK>>
        --
        name : string
        company : string
        email : string <<UK>>
        phone : string
        password : string
        role : UserRole
        is_active : bool
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
        ==
        +ComparePassword(password) : bool
    }

    class Organization {
        +id : uint <<PK>>
        --
        name : string
        phone : string
        #admin_id : uint <<FK>>
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class OrganizationUser {
        +organization_id : uint <<PK,FK>>
        +user_id : uint <<PK,FK>>
    }

    class Language {
        +id : uint <<PK>>
        --
        code : string <<UK>>
        name : string
        description : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class OrganizationLanguage {
        +organization_id : uint <<PK,FK>>
        +language_id : uint <<PK,FK>>
    }

    class Subscription {
        +id : uint <<PK>>
        --
        #organization_id : uint <<FK>>
        period : int
        start_date : timestamp
        end_date : timestamp
        is_active : bool
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class PasswordResetCode {
        +id : uint <<PK>>
        --
        #user_id : uint <<FK>>
        code : string
        expires_at : timestamp
        used : bool
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
        ==
        +IsExpired() : bool
        +IsValid() : bool
    }

    class AdminLog {
        +id : uint <<PK>>
        --
        #admin_id : uint <<FK>>
        action : AdminAction
        entity_type : string
        entity_id : uint
        details : text
        ip_address : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }
}

' =====================================================
' RESTAURANT
' =====================================================
package "Restaurant" <<Frame>> {
    class Restaurant {
        +id : uint <<PK>>
        --
        #organization_id : uint <<FK>>
        name : string
        address : string
        phone : string
        website : string
        description : string
        image_url : string
        menu : string
        reservation_duration : int
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class WorkingHour {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        day_of_week : int
        open_time : string
        close_time : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class Table {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        name : string
        guest_count : int
        menu : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }
}

' =====================================================
' RESERVATION
' =====================================================
package "Reservations" <<Frame>> {
    class Reservation {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        #table_id : uint <<FK>>
        customer_name : string
        customer_phone : string
        customer_email : string
        guest_count : int
        reservation_date : date
        start_time : string
        end_time : string
        status : ReservationStatus
        notes : string
        #chat_session_id : uint <<FK>>
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }
}

' =====================================================
' MENU
' =====================================================
package "Menu" <<Frame>> {
    class MenuCategory {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        name : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class Dish {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        #menu_category_id : uint <<FK>>
        name : string
        price : float
        description : string
        image : string
        is_dish_of_day : bool
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class Ingredient {
        +id : uint <<PK>>
        --
        #dish_id : uint <<FK>>
        name : string
        quantity : float
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class Allergen {
        +id : uint <<PK>>
        --
        #dish_id : uint <<FK>>
        name : string
        description : string
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }
}

' =====================================================
' CHAT
' =====================================================
package "Chat" <<Frame>> {
    class Question {
        +id : uint <<PK>>
        --
        #organization_id : uint <<FK>>
        text : string
        #language_id : uint <<FK>>
        chat_type : ChatType
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class TableChatSession {
        +id : uint <<PK>>
        --
        #table_id : uint <<FK>>
        #restaurant_id : uint <<FK>>
        active : bool
        last_active : timestamp
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class TableChatMessage {
        +id : uint <<PK>>
        --
        #table_id : uint <<FK>>
        #restaurant_id : uint <<FK>>
        #chat_session_id : uint <<FK>>
        author_type : AuthorType
        content : string
        sent_at : timestamp
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class RestaurantChatSession {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        active : bool
        last_active : timestamp
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }

    class RestaurantChatMessage {
        +id : uint <<PK>>
        --
        #restaurant_id : uint <<FK>>
        #chat_session_id : uint <<FK>>
        author_type : AuthorType
        content : string
        sent_at : timestamp
        --
        created_at : timestamp
        updated_at : timestamp
        deleted_at : timestamp
    }
}

' =====================================================
' RELATIONSHIPS
' =====================================================

' User & Organization
User "1" --> "0..*" Organization : administers
Organization "0..*" -- "0..*" User : members
(Organization, User) .. OrganizationUser

' Languages
Organization "0..*" -- "0..*" Language
(Organization, Language) .. OrganizationLanguage

' Subscription
Organization "1" --> "0..1" Subscription

' Password Reset
User "1" --> "0..*" PasswordResetCode

' Admin Logs
User "1" --> "0..*" AdminLog : performs

' Restaurant
Organization "1" --> "0..*" Restaurant
Restaurant "1" --> "0..*" WorkingHour
Restaurant "1" --> "0..*" Table

' Reservations
Restaurant "1" --> "0..*" Reservation
Table "1" --> "0..*" Reservation

' Menu
Restaurant "1" --> "0..*" MenuCategory
MenuCategory "1" --> "0..*" Dish
Restaurant "1" --> "0..*" Dish
Dish "1" --> "0..*" Ingredient
Dish "1" --> "0..*" Allergen

' Questions
Organization "1" --> "0..*" Question
Language "1" --> "0..*" Question

' Table Chat
Table "1" --> "0..*" TableChatSession
Restaurant "1" --> "0..*" TableChatSession
TableChatSession "1" --> "0..*" TableChatMessage

' Restaurant Chat
Restaurant "1" --> "0..*" RestaurantChatSession
RestaurantChatSession "1" --> "0..*" RestaurantChatMessage

@enduml
```

### Как отрендерить PlantUML

**Онлайн:**
- https://www.plantuml.com/plantuml/uml/
- Скопировать код между `@startuml` и `@enduml`

**VS Code:**
- Установить расширение "PlantUML"
- Нажать `Alt+D` для превью

**CLI:**
```bash
# Установка (macOS)
brew install plantuml

# Генерация PNG
plantuml docs/database_schema.md -o output/

# Генерация SVG
plantuml -tsvg docs/database_schema.md -o output/
```

**IntelliJ IDEA / GoLand:**
- Установить плагин "PlantUML Integration"
- Открыть .puml файл или блок в markdown
