package storage

import (
	"gorm.io/gorm"
)

// MenuCategory represents a category in the menu.
// Категория меню привязана к конкретному ресторану.
type MenuCategory struct {
	gorm.Model
	RestaurantID uint       `gorm:"column:restaurant_id;not null;index" json:"restaurant_id"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	Name         string     `gorm:"column:name;not null" json:"name"`
}
