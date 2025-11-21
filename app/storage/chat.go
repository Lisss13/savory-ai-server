package storage

import (
	"gorm.io/gorm"
	"time"
)

type TableChatSessions struct {
	gorm.Model
	TableID      uint                `gorm:"column:table_id;not null" json:"table_id"`
	Table        Table               `gorm:"foreignKey:TableID" json:"table"`
	RestaurantID uint                `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	Restaurant   Restaurant          `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	Active       bool                `gorm:"column:active;not null;default:true" json:"active"`
	LastActive   time.Time           `gorm:"column:last_active;not null" json:"last_active"`
	Messages     []*TableChatMessage `gorm:"foreignKey:ChatSessionID" json:"messages"`
}

type TableChatMessage struct {
	gorm.Model
	TableID          uint              `gorm:"column:table_id;not null" json:"table_id"`
	Table            Table             `gorm:"foreignKey:TableID" json:"table"`
	RestaurantID     uint              `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	Restaurant       Restaurant        `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	ChatSessionID    uint              `gorm:"column:chat_session_id;not null" json:"chat_session_id"`
	TableChatSession TableChatSessions `gorm:"foreignKey:ChatSessionID" json:"table_chat_session"`
	AuthorType       string            `gorm:"column:author_type;not null" json:"author_type"` // "user" || "bot" || "restaurant"
	Content          string            `gorm:"column:content;not null" json:"content"`
	SentAt           time.Time         `gorm:"column:sent_at;not null" json:"sent_at"`
}

type RestaurantChatSessions struct {
	gorm.Model
	RestaurantID uint                     `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	Restaurant   Restaurant               `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	Active       bool                     `gorm:"column:active;not null;default:true" json:"active"`
	LastActive   time.Time                `gorm:"column:last_active;not null" json:"last_active"`
	Messages     []*RestaurantChatMessage `gorm:"foreignKey:ChatSessionID" json:"messages"`
}

type RestaurantChatMessage struct {
	gorm.Model
	RestaurantID          uint                   `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	Restaurant            Restaurant             `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	ChatSessionID         uint                   `gorm:"column:chat_session_id;not null" json:"chat_session_id"`
	RestaurantChatSession RestaurantChatSessions `gorm:"foreignKey:ChatSessionID" json:"restaurant_chat_session"`
	AuthorType            string                 `gorm:"column:author_type;not null" json:"author_type"` // "user" || "bot" || "restaurant"
	Content               string                 `gorm:"column:content;not null" json:"content"`
	SentAt                time.Time              `gorm:"column:sent_at;not null" json:"sent_at"`
}

// Author types for chat messages
const (
	UserAuthor       = "user"
	BotAuthor        = "bot"
	RestaurantAuthor = "restaurant"
)
