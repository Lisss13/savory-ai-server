package storage

import (
	"gorm.io/gorm"
	"time"
)

// PasswordResetCode represents a password reset code
type PasswordResetCode struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Code      string    `gorm:"column:code;not null" json:"code"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	Used      bool      `gorm:"column:used;not null;default:false" json:"used"`
}

// IsExpired checks if the password reset code is expired
func (p *PasswordResetCode) IsExpired() bool {
	return p.ExpiresAt.Before(time.Now())
}

// IsValid checks if the password reset code is valid (not expired and not used)
func (p *PasswordResetCode) IsValid() bool {
	return !p.IsExpired() && !p.Used
}
