package storage

import (
	"gorm.io/gorm"
)

// ChatType определяет тип чата, для которого предназначен вопрос.
// Используется для фильтрации вопросов по контексту использования.
type ChatType string

const (
	// ChatTypeReservation — вопросы для чата бронирования столиков.
	// Пример: "На какое время вы хотите забронировать?", "Сколько гостей?"
	ChatTypeReservation ChatType = "reservation"

	// ChatTypeMenu — вопросы для чата с информацией о меню.
	// Пример: "Какие блюда вы рекомендуете?", "Есть ли вегетарианские опции?"
	ChatTypeMenu ChatType = "menu"
)

// Question представляет вопрос для быстрой отправки в чат-бот.
// Вопросы группируются по организации, языку и типу чата.
type Question struct {
	gorm.Model
	OrganizationID uint         `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Text           string       `gorm:"column:text;not null" json:"text"`
	LanguageID     *uint        `gorm:"column:language_id" json:"language_id"`
	Language       *Language    `gorm:"foreignKey:LanguageID" json:"language"`
	ChatType       ChatType     `gorm:"column:chat_type;type:varchar(20);default:'menu'" json:"chat_type"`
}
