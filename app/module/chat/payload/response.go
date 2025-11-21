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

type RestaurantChatSessionsResp struct {
	Sessions []TableChatSessionsResp `json:"sessions"`
}
