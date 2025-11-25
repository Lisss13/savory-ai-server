package storage

import (
	"time"

	"gorm.io/gorm"
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
	ReservationStatusNoShow    ReservationStatus = "no_show"
)

// Reservation represents a table reservation
type Reservation struct {
	gorm.Model
	RestaurantID    uint              `gorm:"column:restaurant_id;not null;index" json:"restaurant_id"`
	Restaurant      Restaurant        `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	TableID         uint              `gorm:"column:table_id;not null;index" json:"table_id"`
	Table           Table             `gorm:"foreignKey:TableID" json:"table"`
	CustomerName    string            `gorm:"column:customer_name;not null" json:"customer_name"`
	CustomerPhone   string            `gorm:"column:customer_phone;not null" json:"customer_phone"`
	CustomerEmail   string            `gorm:"column:customer_email" json:"customer_email"`
	GuestCount      int               `gorm:"column:guest_count;not null" json:"guest_count"`
	ReservationDate time.Time         `gorm:"column:reservation_date;not null;index" json:"reservation_date"`
	StartTime       string            `gorm:"column:start_time;not null" json:"start_time"` // Format: "HH:MM"
	EndTime         string            `gorm:"column:end_time;not null" json:"end_time"`     // Format: "HH:MM"
	Status          ReservationStatus `gorm:"column:status;not null;default:'confirmed'" json:"status"`
	Notes           string            `gorm:"column:notes" json:"notes"`
	ChatSessionID   *uint             `gorm:"column:chat_session_id" json:"chat_session_id"` // Reference to a chat session if booked via chat
}
