package repository

import (
	"errors"
	"gorm.io/gorm"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type dishRepository struct {
	DB *database.Database
}

type DishRepository interface {
	FindAll() (dishes []*storage.Dish, err error)
	FindByID(id uint) (dish *storage.Dish, err error)
	FindByOrganizationID(organizationID uint) (dishes []*storage.Dish, err error)
	Create(dish *storage.Dish) (res *storage.Dish, err error)
	Update(dish *storage.Dish) (res *storage.Dish, err error)
	Delete(id uint) error
	FindDishOfDay(companyID uint) (dish *storage.Dish, err error)
	SetDishOfDay(id uint) (dish *storage.Dish, err error)
}

func NewDishRepository(db *database.Database) DishRepository {
	return &dishRepository{
		DB: db,
	}
}

func (r *dishRepository) FindAll() (dishes []*storage.Dish, err error) {
	if err := r.DB.DB.Preload("Organization").Preload("MenuCategory").Preload("Ingredients").Preload("Allergens").Find(&dishes).Error; err != nil {
		return nil, err
	}
	return dishes, nil
}

func (r *dishRepository) FindByOrganizationID(organizationID uint) (dishes []*storage.Dish, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload("MenuCategory").
		Preload("Ingredients").
		Preload("Allergens").
		Where("organization_id = ?", organizationID).
		Find(&dishes).Error

	if err != nil {
		return nil, err
	}

	return dishes, nil
}

func (r *dishRepository) FindByID(id uint) (dish *storage.Dish, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload("MenuCategory").
		Preload("Ingredients").
		Preload("Allergens").
		First(&dish, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return dish, nil
}

func (r *dishRepository) Create(dish *storage.Dish) (res *storage.Dish, err error) {
	if err := r.DB.DB.Create(&dish).Error; err != nil {
		return nil, err
	}

	// Reload the dish with all associations
	return r.FindByID(dish.ID)
}

func (r *dishRepository) Update(dish *storage.Dish) (res *storage.Dish, err error) {
	// First, update the dish itself
	if err := r.DB.DB.Model(&dish).Updates(map[string]interface{}{
		"organization_id":  dish.OrganizationID,
		"menu_category_id": dish.MenuCategoryID,
		"name":             dish.Name,
		"price":            dish.Price,
		"description":      dish.Description,
		"image":            dish.Image,
	}).Error; err != nil {
		return nil, err
	}

	// Handle ingredients update
	// Delete existing ingredients
	if err = r.DB.DB.Where("dish_id = ?", dish.ID).Delete(&storage.Ingredient{}).Error; err != nil {
		return nil, err
	}

	// Create new ingredients if provided
	if len(dish.Ingredients) > 0 {
		for _, ingredient := range dish.Ingredients {
			ingredient.ID = 0 // Ensure new ingredients are created
			ingredient.DishID = dish.ID
		}
		if err = r.DB.DB.Create(&dish.Ingredients).Error; err != nil {
			return nil, err
		}
	}

	// Handle allergens update
	// Delete existing allergens
	if err := r.DB.DB.Where("dish_id = ?", dish.ID).Delete(&storage.Allergen{}).Error; err != nil {
		return nil, err
	}

	// Create new allergens if provided
	if len(dish.Allergens) > 0 {
		for _, allergen := range dish.Allergens {
			allergen.ID = 0 // Ensure new allergens are created
			allergen.DishID = dish.ID
		}
		if err := r.DB.DB.Create(&dish.Allergens).Error; err != nil {
			return nil, err
		}
	}

	// Reload the dish with all associations
	return r.FindByID(dish.ID)
}

func (r *dishRepository) Delete(id uint) error {
	// First delete all ingredients
	if err := r.DB.DB.Where("dish_id = ?", id).Delete(&storage.Ingredient{}).Error; err != nil {
		return err
	}

	// Delete all allergens
	if err := r.DB.DB.Where("dish_id = ?", id).Delete(&storage.Allergen{}).Error; err != nil {
		return err
	}

	// Then delete the dish
	return r.DB.DB.Delete(&storage.Dish{}, id).Error
}

func (r *dishRepository) FindDishOfDay(companyID uint) (dish *storage.Dish, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload("MenuCategory").
		Preload("Ingredients").
		Preload("Allergens").
		Where("is_dish_of_day = ?", true).
		Where("organization_id = ?", companyID).
		First(&dish).
		Error

	if err != nil {
		return nil, err
	}

	return dish, nil
}

func (r *dishRepository) SetDishOfDay(id uint) (dish *storage.Dish, err error) {
	// First, reset all dishes to not be dish of the day
	err = r.DB.DB.
		Model(&storage.Dish{}).
		Where("is_dish_of_day = ?", true).
		Update("is_dish_of_day", false).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Then, set the specified dish as dish of the day
	err = r.DB.DB.
		Model(&storage.Dish{}).
		Where("id = ?", id).
		Update("is_dish_of_day", true).
		Error

	if err != nil {
		return nil, err
	}

	// Return the updated dish
	return r.FindByID(id)
}
