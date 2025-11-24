package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type adminRepository struct {
	DB *database.Database
}

type AdminRepository interface {
	// Statistics
	CountUsers() (int64, error)
	CountActiveUsers() (int64, error)
	CountOrganizations() (int64, error)
	CountRestaurants() (int64, error)
	CountDishes() (int64, error)
	CountTables() (int64, error)
	CountQuestions() (int64, error)
	CountActiveSubscriptions() (int64, error)

	// Users
	FindAllUsers(page, pageSize int) ([]*storage.User, int64, error)
	FindUserByID(id uint) (*storage.User, error)
	UpdateUserStatus(id uint, isActive bool) error
	UpdateUserRole(id uint, role storage.UserRole) error
	DeleteUser(id uint) error

	// Organizations
	FindAllOrganizations(page, pageSize int) ([]*storage.Organization, int64, error)
	FindOrganizationByID(id uint) (*storage.Organization, error)
	DeleteOrganization(id uint) error

	// Dishes (for moderation)
	FindAllDishes(page, pageSize int) ([]*storage.Dish, int64, error)
	DeleteDish(id uint) error

	// Admin Logs
	CreateLog(log *storage.AdminLog) error
	FindAllLogs(page, pageSize int) ([]*storage.AdminLog, int64, error)
	FindLogsByAdminID(adminID uint, page, pageSize int) ([]*storage.AdminLog, int64, error)
}

func NewAdminRepository(db *database.Database) AdminRepository {
	return &adminRepository{DB: db}
}

// ==================== Statistics ====================

func (r *adminRepository) CountUsers() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.User{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountActiveUsers() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.User{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountOrganizations() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Organization{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountRestaurants() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Restaurant{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountDishes() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Dish{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountTables() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Table{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountQuestions() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Question{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) CountActiveSubscriptions() (int64, error) {
	var count int64
	err := r.DB.DB.Model(&storage.Subscription{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

// ==================== Users ====================

func (r *adminRepository) FindAllUsers(page, pageSize int) ([]*storage.User, int64, error) {
	var users []*storage.User
	var count int64

	offset := (page - 1) * pageSize

	if err := r.DB.DB.Model(&storage.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *adminRepository) FindUserByID(id uint) (*storage.User, error) {
	var user storage.User
	if err := r.DB.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *adminRepository) UpdateUserStatus(id uint, isActive bool) error {
	return r.DB.DB.Model(&storage.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *adminRepository) UpdateUserRole(id uint, role storage.UserRole) error {
	return r.DB.DB.Model(&storage.User{}).Where("id = ?", id).Update("role", role).Error
}

func (r *adminRepository) DeleteUser(id uint) error {
	return r.DB.DB.Delete(&storage.User{}, id).Error
}

// ==================== Organizations ====================

func (r *adminRepository) FindAllOrganizations(page, pageSize int) ([]*storage.Organization, int64, error) {
	var organizations []*storage.Organization
	var count int64

	offset := (page - 1) * pageSize

	if err := r.DB.DB.Model(&storage.Organization{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.DB.Preload("Admin").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&organizations).Error; err != nil {
		return nil, 0, err
	}

	return organizations, count, nil
}

func (r *adminRepository) FindOrganizationByID(id uint) (*storage.Organization, error) {
	var organization storage.Organization
	if err := r.DB.DB.Preload("Admin").First(&organization, id).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *adminRepository) DeleteOrganization(id uint) error {
	return r.DB.DB.Delete(&storage.Organization{}, id).Error
}

// ==================== Dishes ====================

func (r *adminRepository) FindAllDishes(page, pageSize int) ([]*storage.Dish, int64, error) {
	var dishes []*storage.Dish
	var count int64

	offset := (page - 1) * pageSize

	if err := r.DB.DB.Model(&storage.Dish{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.DB.Preload("Organization").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&dishes).Error; err != nil {
		return nil, 0, err
	}

	return dishes, count, nil
}

func (r *adminRepository) DeleteDish(id uint) error {
	// Delete ingredients first
	if err := r.DB.DB.Where("dish_id = ?", id).Delete(&storage.Ingredient{}).Error; err != nil {
		return err
	}
	// Delete allergens
	if err := r.DB.DB.Where("dish_id = ?", id).Delete(&storage.Allergen{}).Error; err != nil {
		return err
	}
	// Delete dish
	return r.DB.DB.Delete(&storage.Dish{}, id).Error
}

// ==================== Admin Logs ====================

func (r *adminRepository) CreateLog(log *storage.AdminLog) error {
	return r.DB.DB.Create(log).Error
}

func (r *adminRepository) FindAllLogs(page, pageSize int) ([]*storage.AdminLog, int64, error) {
	var logs []*storage.AdminLog
	var count int64

	offset := (page - 1) * pageSize

	if err := r.DB.DB.Model(&storage.AdminLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.DB.Preload("Admin").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}

func (r *adminRepository) FindLogsByAdminID(adminID uint, page, pageSize int) ([]*storage.AdminLog, int64, error) {
	var logs []*storage.AdminLog
	var count int64

	offset := (page - 1) * pageSize

	if err := r.DB.DB.Model(&storage.AdminLog{}).Where("admin_id = ?", adminID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.DB.Preload("Admin").Where("admin_id = ?", adminID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}
