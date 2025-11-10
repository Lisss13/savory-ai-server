package storage

import (
	"gorm.io/gorm"
)

// MenuCategory represents a category in the menu
type MenuCategory struct {
	gorm.Model
	OrganizationID uint         `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Name           string       `gorm:"column:name;not null" json:"name"`
}
