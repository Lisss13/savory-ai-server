// Package payload содержит структуры запросов и ответов для модуля чата.
package payload

// =====================================================
// Table Chat Requests - запросы для чата столика
// =====================================================

// StartTableSessionReq запрос на создание сессии чата для столика.
// Используется посетителем при сканировании QR-кода.
type StartTableSessionReq struct {
	TableID      uint `json:"tableId" validate:"required"`      // ID столика
	RestaurantID uint `json:"restaurantId" validate:"required"` // ID ресторана
}

// CloseTableSessionReq запрос на закрытие сессии чата столика.
type CloseTableSessionReq struct {
	SessionID uint `json:"sessionId" validate:"required"` // ID сессии для закрытия
}

// SendTableMessageReq запрос на отправку сообщения в чат столика.
type SendTableMessageReq struct {
	SessionID uint   `json:"sessionId" validate:"required"` // ID активной сессии
	Content   string `json:"content" validate:"required"`   // Текст сообщения
}

// =====================================================
// Restaurant Chat Requests - запросы для чата ресторана
// =====================================================

// StartRestaurantSessionReq запрос на создание сессии чата с рестораном.
// Используется для общения с AI-ботом (бронирование, вопросы).
type StartRestaurantSessionReq struct {
	RestaurantID uint `json:"restaurantId" validate:"required"` // ID ресторана
}

// CloseRestaurantSessionReq запрос на закрытие сессии чата ресторана.
type CloseRestaurantSessionReq struct {
	SessionID uint `json:"sessionId" validate:"required"` // ID сессии для закрытия
}

// SendRestaurantMessageReq запрос на отправку сообщения в чат ресторана.
// AI обрабатывает сообщение и может выполнять tool calls для бронирования.
type SendRestaurantMessageReq struct {
	SessionID uint   `json:"sessionId" validate:"required"` // ID активной сессии
	Content   string `json:"content" validate:"required"`   // Текст сообщения
}
