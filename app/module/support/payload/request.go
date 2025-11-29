// Package payload содержит структуры запросов и ответов для модуля поддержки.
package payload

// CreateSupportTicketReq запрос на создание заявки в поддержку
type CreateSupportTicketReq struct {
	Title       string `json:"title" validate:"required"`       // Заголовок проблемы
	Description string `json:"description" validate:"required"` // Описание проблемы
	Email       string `json:"email" validate:"required,email"` // Email для связи (обязательно)
	Phone       string `json:"phone" validate:"omitempty"`      // Телефон для связи (опционально)
}

// UpdateSupportTicketStatusReq запрос на обновление статуса заявки (только для admin)
type UpdateSupportTicketStatusReq struct {
	Status string `json:"status" validate:"required,oneof=in_progress completed"` // Статус: in_progress или completed
}
