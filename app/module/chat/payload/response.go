// Package payload содержит структуры запросов и ответов для модуля чата.
package payload

import (
	"time"
)

// =====================================================
// Table Chat Responses - ответы для чата столика
// =====================================================

// TableResp краткая информация о столике.
type TableResp struct {
	ID   uint   `json:"id"`   // ID столика
	Name string `json:"name"` // Название столика
}

// TableMessageResp сообщение в чате столика.
type TableMessageResp struct {
	ID      uint      `json:"id"`      // ID сообщения
	Content string    `json:"content"` // Текст сообщения
	SentAt  time.Time `json:"sentAt"`  // Время отправки
}

// TableChatSessions данные сессии чата столика.
type TableChatSessions struct {
	ID         uint               `json:"id"`         // ID сессии
	Active     bool               `json:"active"`     // Активна ли сессия
	LastActive time.Time          `json:"lastActive"` // Время последней активности
	Table      TableResp          `json:"table"`      // Информация о столике
	Messages   []TableMessageResp `json:"messages"`   // Сообщения сессии
}

// TableChatSessionsResp обёртка для ответа с сессией чата столика.
type TableChatSessionsResp struct {
	Session TableChatSessions `json:"session"` // Данные сессии
}

// Message универсальная структура сообщения.
type Message struct {
	ID      uint      `json:"id"`      // ID сообщения
	Content string    `json:"content"` // Текст сообщения
	SentAt  time.Time `json:"sentAt"`  // Время отправки
}

// MessagesRespFormBot ответ с сообщением от бота (чат столика).
type MessagesRespFormBot struct {
	Message Message `json:"message"` // Сообщение от AI-бота
}

// TableChatMessageResp сообщение с указанием автора.
type TableChatMessageResp struct {
	ID         uint      `json:"id"`         // ID сообщения
	Content    string    `json:"content"`    // Текст сообщения
	SentAt     time.Time `json:"sentAt"`     // Время отправки
	AuthorType string    `json:"authorType"` // Тип автора: "user" или "bot"
}

// TableChatMessagesResp список сообщений сессии чата столика.
type TableChatMessagesResp struct {
	Messages []TableChatMessageResp `json:"messages"` // Сообщения сессии
}

// TableChatSessionsByTableIDResp список сессий для столика.
type TableChatSessionsByTableIDResp struct {
	Sessions []TableChatSessionsResp `json:"sessions"` // Сессии столика
}

// =====================================================
// Restaurant Chat Responses - ответы для чата ресторана
// =====================================================

// RestaurantResp краткая информация о ресторане.
type RestaurantResp struct {
	ID   uint   `json:"id"`   // ID ресторана
	Name string `json:"name"` // Название ресторана
}

// RestaurantMessageResp сообщение в чате ресторана.
type RestaurantMessageResp struct {
	ID      uint      `json:"id"`      // ID сообщения
	Content string    `json:"content"` // Текст сообщения
	SentAt  time.Time `json:"sentAt"`  // Время отправки
}

// RestaurantChatSession данные сессии чата ресторана.
type RestaurantChatSession struct {
	ID         uint                    `json:"id"`         // ID сессии
	Active     bool                    `json:"active"`     // Активна ли сессия
	LastActive time.Time               `json:"lastActive"` // Время последней активности
	Restaurant RestaurantResp          `json:"restaurant"` // Информация о ресторане
	Messages   []RestaurantMessageResp `json:"messages"`   // Сообщения сессии
}

// RestaurantChatSessionResp обёртка для ответа с сессией чата ресторана.
type RestaurantChatSessionResp struct {
	Session RestaurantChatSession `json:"session"` // Данные сессии
}

// RestaurantChatSessionsResp список сессий чата ресторана.
type RestaurantChatSessionsResp struct {
	Sessions []RestaurantChatSessionResp `json:"sessions"` // Сессии ресторана
}

// RestaurantChatMessageResp сообщение чата ресторана с указанием автора.
type RestaurantChatMessageResp struct {
	ID         uint      `json:"id"`         // ID сообщения
	Content    string    `json:"content"`    // Текст сообщения
	SentAt     time.Time `json:"sentAt"`     // Время отправки
	AuthorType string    `json:"authorType"` // Тип автора: "user" или "bot"
}

// RestaurantChatMessagesResp список сообщений сессии чата ресторана.
type RestaurantChatMessagesResp struct {
	Messages []RestaurantChatMessageResp `json:"messages"` // Сообщения сессии
}

// RestaurantMessagesRespFormBot ответ с сообщением от бота (чат ресторана).
type RestaurantMessagesRespFormBot struct {
	Message Message `json:"message"` // Сообщение от AI-бота
}
