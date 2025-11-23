package storage

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name      string     `gorm:"column:name;not null" json:"name"`
	Phone     string     `gorm:"column:phone;not null" json:"phone"`
	AdminID   uint       `gorm:"column:admin_id;not null" json:"admin_id"`
	Admin     User       `gorm:"foreignKey:AdminID" json:"admin"`
	Users     []User     `gorm:"many2many:organization_users;" json:"users"`
	Languages []Language `gorm:"many2many:organization_languages;" json:"languages"`
}

// OrganizationUser represents the many-to-many relationship between organizations and users
type OrganizationUser struct {
	OrganizationID uint `gorm:"primaryKey"`
	UserID         uint `gorm:"primaryKey"`
}
