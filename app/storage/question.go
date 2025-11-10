package storage

import (
	"gorm.io/gorm"
)

// Question represents a question
type Question struct {
	gorm.Model
	OrganizationID uint         `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Text           string       `gorm:"column:text;not null" json:"text"`
}
