package storage

import (
	"gorm.io/gorm"
)

// Restaurant represents a restaurant
type Restaurant struct {
	gorm.Model
	OrganizationID uint           `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID" json:"organization"`
	Name           string         `gorm:"column:name;not null" json:"name"`
	Address        string         `gorm:"column:address;not null" json:"address"`
	Phone          string         `gorm:"column:phone;not null" json:"phone"`
	Website        string         `gorm:"column:website" json:"website"`
	Description    string         `gorm:"column:description" json:"description"`
	ImageURL       string         `gorm:"column:image_url" json:"image_url"`
	Menu           string         `gorm:"column:menu" json:"menu"`
	WorkingHours   []*WorkingHour `gorm:"foreignKey:RestaurantID" json:"working_hours"`
}

// WorkingHour represents a restaurant's working hours for a specific day
type WorkingHour struct {
	gorm.Model
	RestaurantID uint   `gorm:"column:restaurant_id;not null" json:"restaurant_id"`
	DayOfWeek    int    `gorm:"column:day_of_week;not null" json:"day_of_week"` // 0 = Sunday, 1 = Monday, etc.
	OpenTime     string `gorm:"column:open_time;not null" json:"open_time"`
	CloseTime    string `gorm:"column:close_time;not null" json:"close_time"`
}
