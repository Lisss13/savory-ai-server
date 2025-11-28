package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type menuCategoryRepository struct {
	DB *database.Database
}

// MenuCategoryRepository определяет интерфейс для работы с категориями меню в БД.
type MenuCategoryRepository interface {
	FindAll() (categories []*storage.MenuCategory, err error)
	FindByID(id uint) (category *storage.MenuCategory, err error)
	FindByRestaurantID(restaurantID uint) (categories []*storage.MenuCategory, err error)
	Create(category *storage.MenuCategory) (res *storage.MenuCategory, err error)
	Delete(id uint) error
}

func NewMenuCategoryRepository(db *database.Database) MenuCategoryRepository {
	return &menuCategoryRepository{
		DB: db,
	}
}

func (r *menuCategoryRepository) FindAll() (categories []*storage.MenuCategory, err error) {
	if err := r.DB.DB.Preload("Restaurant").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByRestaurantID возвращает все категории меню для указанного ресторана.
func (r *menuCategoryRepository) FindByRestaurantID(restaurantID uint) (categories []*storage.MenuCategory, err error) {
	if err := r.DB.DB.Preload("Restaurant").Where("restaurant_id = ?", restaurantID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *menuCategoryRepository) FindByID(id uint) (category *storage.MenuCategory, err error) {
	err = r.DB.DB.
		Preload("Restaurant").
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
