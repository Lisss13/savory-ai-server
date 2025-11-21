package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type userRepository struct {
	DB *database.Database
}

type UserRepository interface {
	FindUserByEmail(email string) (user *storage.User, err error)
	CreateUser(user *storage.User) (res *storage.User, err error)
	FindUserByID(id int64) (user *storage.User, err error)
	UpdatePassword(userID uint, password string) error
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
