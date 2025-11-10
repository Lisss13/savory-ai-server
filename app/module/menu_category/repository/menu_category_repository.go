package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type menuCategoryRepository struct {
	DB *database.Database
}

type MenuCategoryRepository interface {
	FindAll() (categories []*storage.MenuCategory, err error)
	FindByID(id uint) (category *storage.MenuCategory, err error)
	FindByOrganizationID(organizationID uint) (categories []*storage.MenuCategory, err error)
	Create(category *storage.MenuCategory) (res *storage.MenuCategory, err error)
	Delete(id uint) error
}

func NewMenuCategoryRepository(db *database.Database) MenuCategoryRepository {
	return &menuCategoryRepository{
		DB: db,
	}
}

func (r *menuCategoryRepository) FindAll() (categories []*storage.MenuCategory, err error) {
	if err := r.DB.DB.Preload("Organization").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *menuCategoryRepository) FindByOrganizationID(organizationID uint) (categories []*storage.MenuCategory, err error) {
	if err := r.DB.DB.Preload("Organization").Where("organization_id = ?", organizationID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *menuCategoryRepository) FindByID(id uint) (category *storage.MenuCategory, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload(clause.Associations).
		First(&category, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *menuCategoryRepository) Create(category *storage.MenuCategory) (res *storage.MenuCategory, err error) {
	if err := r.DB.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (r *menuCategoryRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.MenuCategory{}, id).Error
}
