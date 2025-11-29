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
	// =====================================================
	// CRUD Operations
	// =====================================================
	FindAll() (categories []*storage.MenuCategory, err error)
	FindByID(id uint) (category *storage.MenuCategory, err error)
	FindByRestaurantID(restaurantID uint) (categories []*storage.MenuCategory, err error)
	Create(category *storage.MenuCategory) (res *storage.MenuCategory, err error)
	Update(category *storage.MenuCategory) (res *storage.MenuCategory, err error)
	Delete(id uint) error

	// =====================================================
	// Sort Order Operations
	// =====================================================
	// UpdateSortOrder обновляет порядок сортировки для категории.
	UpdateSortOrder(id uint, sortOrder int) error
	// UpdateSortOrderBatch массово обновляет порядок сортировки для нескольких категорий.
	UpdateSortOrderBatch(updates map[uint]int) error
	// GetMaxSortOrder возвращает максимальный sort_order для ресторана.
	GetMaxSortOrder(restaurantID uint) (int, error)
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
// Категории сортируются по полю sort_order (по возрастанию).
func (r *menuCategoryRepository) FindByRestaurantID(restaurantID uint) (categories []*storage.MenuCategory, err error) {
	if err := r.DB.DB.Preload("Restaurant").
		Where("restaurant_id = ?", restaurantID).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error; err != nil {
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

// Update обновляет категорию меню.
func (r *menuCategoryRepository) Update(category *storage.MenuCategory) (res *storage.MenuCategory, err error) {
	if err := r.DB.DB.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateSortOrder обновляет порядок сортировки для конкретной категории.
func (r *menuCategoryRepository) UpdateSortOrder(id uint, sortOrder int) error {
	return r.DB.DB.Model(&storage.MenuCategory{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

// UpdateSortOrderBatch массово обновляет порядок сортировки для нескольких категорий.
// Принимает map[categoryID]sortOrder.
func (r *menuCategoryRepository) UpdateSortOrderBatch(updates map[uint]int) error {
	tx := r.DB.DB.Begin()
	for id, sortOrder := range updates {
		if err := tx.Model(&storage.MenuCategory{}).
			Where("id = ?", id).
			Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// GetMaxSortOrder возвращает максимальный sort_order для ресторана.
// Используется при создании новой категории для автоматической установки порядка.
func (r *menuCategoryRepository) GetMaxSortOrder(restaurantID uint) (int, error) {
	var maxSortOrder int
	err := r.DB.DB.Model(&storage.MenuCategory{}).
		Where("restaurant_id = ?", restaurantID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error
	return maxSortOrder, err
}
