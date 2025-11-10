package storage

import (
	"gorm.io/gorm"
)

// Table represents a restaurant table
type Table struct {
	gorm.Model
	RestaurantID uint       `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	Name         string     `gorm:"column:name;not null" json:"name"`
	GuestCount   int        `gorm:"column:guest_count;not null" json:"guest_count"`
	QRCodeURL    string     `gorm:"column:qr_code_url" json:"qr_code_url"`
}
