package storage

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;not null" json:"name"`
	Company  string `gorm:"column:company;not null" json:"company"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Phone    string `gorm:"column:phone;not null" json:"phone"`
	Password string `gorm:"column:password;not null" json:"-"`
}

// ComparePassword compare password
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
