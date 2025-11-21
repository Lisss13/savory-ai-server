package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
	"time"
)

type userRepository struct {
	DB *database.Database
}

type UserRepository interface {
	FindUserByEmail(email string) (user *storage.User, err error)
	CreateUser(user *storage.User) (res *storage.User, err error)
	FindUserByID(id int64) (user *storage.User, err error)
	UpdatePassword(userID uint, password string) error

	// Password reset methods
	CreatePasswordResetCode(userID uint, code string, expiresAt time.Time) (*storage.PasswordResetCode, error)
	FindPasswordResetCodeByCode(code string) (*storage.PasswordResetCode, error)
	FindPasswordResetCodeByUserEmail(email string) (*storage.PasswordResetCode, error)
	MarkPasswordResetCodeAsUsed(codeID uint) error
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur *userRepository) FindUserByEmail(email string) (user *storage.User, err error) {
	if err := ur.DB.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) CreateUser(user *storage.User) (res *storage.User, err error) {
	if err := ur.DB.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) FindUserByID(id int64) (user *storage.User, err error) {
	err = ur.DB.DB.
		Preload(clause.Associations).
		First(&user, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) UpdatePassword(userID uint, password string) error {
	return ur.DB.DB.Model(&storage.User{}).
		Where("id = ?", userID).
		Update("password", password).
		Error
}

// CreatePasswordResetCode creates a new password reset code
func (ur *userRepository) CreatePasswordResetCode(userID uint, code string, expiresAt time.Time) (*storage.PasswordResetCode, error) {
	// Create a new password reset code
	passwordResetCode := &storage.PasswordResetCode{
		UserID:    userID,
		Code:      code,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	// Save to database
	if err := ur.DB.DB.Create(&passwordResetCode).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	var createdCode storage.PasswordResetCode
	if err := ur.DB.DB.
		Preload("User").
		First(&createdCode, passwordResetCode.ID).Error; err != nil {
		return nil, err
	}

	return &createdCode, nil
}

// FindPasswordResetCodeByCode finds a password reset code by its code
func (ur *userRepository) FindPasswordResetCodeByCode(code string) (*storage.PasswordResetCode, error) {
	var passwordResetCode storage.PasswordResetCode

	err := ur.DB.DB.
		Preload("User").
		Where("code = ?", code).
		First(&passwordResetCode).Error

	if err != nil {
		return nil, err
	}

	return &passwordResetCode, nil
}

// FindPasswordResetCodeByUserEmail finds a password reset code by user email
func (ur *userRepository) FindPasswordResetCodeByUserEmail(email string) (*storage.PasswordResetCode, error) {
	// First find the user
	user, err := ur.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Then find the password reset code
	var passwordResetCode storage.PasswordResetCode

	err = ur.DB.DB.
		Preload("User").
		Where("user_id = ? AND used = ? AND expires_at > ?", user.ID, false, time.Now()).
		Order("created_at DESC").
		First(&passwordResetCode).Error

	if err != nil {
		return nil, err
	}

	return &passwordResetCode, nil
}

// MarkPasswordResetCodeAsUsed marks a password reset code as used
func (ur *userRepository) MarkPasswordResetCodeAsUsed(codeID uint) error {
	return ur.DB.DB.Model(&storage.PasswordResetCode{}).
		Where("id = ?", codeID).
		Update("used", true).
		Error
}
