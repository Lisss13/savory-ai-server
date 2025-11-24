package storage

import (
	"gorm.io/gorm"
)

// AdminAction represents types of admin actions
type AdminAction string

const (
	ActionCreate     AdminAction = "create"
	ActionUpdate     AdminAction = "update"
	ActionDelete     AdminAction = "delete"
	ActionBlock      AdminAction = "block"
	ActionUnblock    AdminAction = "unblock"
	ActionActivate   AdminAction = "activate"
	ActionDeactivate AdminAction = "deactivate"
	ActionView       AdminAction = "view"
)

// AdminLog represents an admin action log entry
type AdminLog struct {
	gorm.Model
	AdminID    uint        `gorm:"column:admin_id;not null" json:"admin_id"`
	Admin      User        `gorm:"foreignKey:AdminID" json:"admin"`
	Action     AdminAction `gorm:"column:action;not null" json:"action"`
	EntityType string      `gorm:"column:entity_type;not null" json:"entity_type"` // user, organization, subscription, dish, etc.
	EntityID   uint        `gorm:"column:entity_id" json:"entity_id"`
	Details    string      `gorm:"column:details;type:text" json:"details"` // JSON with additional info
	IPAddress  string      `gorm:"column:ip_address" json:"ip_address"`
}
