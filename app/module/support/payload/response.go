// Package payload содержит структуры запросов и ответов для модуля поддержки.
package payload

import "time"

// SupportTicketResp ответ с данными заявки в поддержку
type SupportTicketResp struct {
	ID          uint      `json:"id"`          // ID заявки
	UserID      uint      `json:"user_id"`     // ID пользователя
	UserName    string    `json:"user_name"`   // Имя пользователя
	UserEmail   string    `json:"user_email"`  // Email пользователя из профиля
	Title       string    `json:"title"`       // Заголовок проблемы
	Description string    `json:"description"` // Описание проблемы
	Email       string    `json:"email"`       // Email для связи (из заявки)
	Phone       string    `json:"phone"`       // Телефон для связи
	Status      string    `json:"status"`      // Статус заявки
	CreatedAt   time.Time `json:"created_at"`  // Дата создания
	UpdatedAt   time.Time `json:"updated_at"`  // Дата обновления
}

// SupportTicketsListResp список заявок с пагинацией
type SupportTicketsListResp struct {
	Tickets    []SupportTicketResp `json:"tickets"`     // Список заявок
	TotalCount int64               `json:"total_count"` // Общее количество
	Page       int                 `json:"page"`        // Текущая страница
	PageSize   int                 `json:"page_size"`   // Размер страницы
}
