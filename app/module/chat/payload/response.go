package payload

import (
	"time"
)

type TableResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TableMessageResp struct {
	ID      uint      `json:"id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sentAt"`
}

type TableChatSessions struct {
	ID         uint               `json:"id"`
	Active     bool               `json:"active"`
	LastActive time.Time          `json:"lastActive"`
	Table      TableResp          `json:"table"`
	Messages   []TableMessageResp `json:"messages"`
}

type TableChatSessionsResp struct {
	Session TableChatSessions `json:"session"`
}

type Message struct {
	ID      uint      `json:"id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sentAt"`
}

type MessagesRespFormBot struct {
	Message Message `json:"message"`
}

type TableChatMessageResp struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sentAt"`
	AuthorType string    `json:"authorType"`
}

type TableChatMessagesResp struct {
	Messages []TableChatMessageResp `json:"messages"`
}

type TableChatSessionsByTableIDResp struct {
	Sessions []TableChatSessionsResp `json:"sessions"`
}

type RestaurantResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RestaurantMessageResp struct {
	ID      uint      `json:"id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sentAt"`
}

type RestaurantChatSession struct {
	ID         uint                  `json:"id"`
	Active     bool                  `json:"active"`
	LastActive time.Time             `json:"lastActive"`
	Restaurant RestaurantResp        `json:"restaurant"`
	Messages   []RestaurantMessageResp `json:"messages"`
}

type RestaurantChatSessionResp struct {
	Session RestaurantChatSession `json:"session"`
}

type RestaurantChatSessionsResp struct {
	Sessions []RestaurantChatSessionResp `json:"sessions"`
}

type RestaurantChatMessageResp struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sentAt"`
	AuthorType string    `json:"authorType"`
}

type RestaurantChatMessagesResp struct {
	Messages []RestaurantChatMessageResp `json:"messages"`
}

type RestaurantMessagesRespFormBot struct {
	Message Message `json:"message"`
}
