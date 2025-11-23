package storage

import (
	"gorm.io/gorm"
	"time"
)

// Subscription represents a subscription for an organization
type Subscription struct {
	gorm.Model
	OrganizationID uint         `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Period         int          `gorm:"column:period;not null" json:"period"`           // Number of months
	StartDate      time.Time    `gorm:"column:start_date;not null" json:"start_date"`   // Subscription start date
	EndDate        time.Time    `gorm:"column:end_date;not null" json:"end_date"`       // Subscription end date
	IsActive       bool         `gorm:"column:is_active;default:true" json:"is_active"` // Is subscription active
}
