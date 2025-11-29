package storage

import (
	"gorm.io/gorm"
)

// SupportTicketStatus represents the status of a support ticket
type SupportTicketStatus string

const (
	SupportTicketStatusInProgress SupportTicketStatus = "in_progress" // Взят в работу
	SupportTicketStatusCompleted  SupportTicketStatus = "completed"   // Завершен
)

// SupportTicket represents a support request from a user
// Заявка в службу поддержки от пользователя системы
type SupportTicket struct {
	gorm.Model
	UserID      uint                `gorm:"column:user_id;not null;index" json:"user_id"`        // ID пользователя, создавшего заявку
	User        User                `gorm:"foreignKey:UserID" json:"user"`                       // Пользователь
	Title       string              `gorm:"column:title;not null" json:"title"`                  // Заголовок проблемы
	Description string              `gorm:"column:description;type:text;not null" json:"description"` // Описание проблемы
	Email       string              `gorm:"column:email;not null" json:"email"`                  // Email для связи (обязательно)
	Phone       string              `gorm:"column:phone" json:"phone"`                           // Телефон для связи (опционально)
	Status      SupportTicketStatus `gorm:"column:status;not null;default:'in_progress'" json:"status"` // Статус заявки
}
