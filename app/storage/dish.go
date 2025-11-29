package storage

import (
	"gorm.io/gorm"
)

// Dish represents a dish in the menu.
// Блюдо привязано к ресторану и категории меню.
type Dish struct {
	gorm.Model
	RestaurantID   uint          `gorm:"column:restaurant_id;not null;index" json:"restaurant_id"`
	Restaurant     Restaurant    `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	MenuCategoryID uint          `gorm:"column:menu_category_id;not null" json:"menu_category_id"`
	MenuCategory   MenuCategory  `gorm:"foreignKey:MenuCategoryID" json:"menu_category"`
	Name           string        `gorm:"column:name;not null" json:"name"`
	Price          float64       `gorm:"column:price;not null" json:"price"`
	Description    string        `gorm:"column:description" json:"description"`
	Image          string        `gorm:"column:image" json:"image"`
	IsDishOfDay    bool          `gorm:"column:is_dish_of_day;default:false" json:"is_dish_of_day"`
	Proteins       float64       `gorm:"column:proteins;default:0" json:"proteins"`           // Белки (г)
	Fats           float64       `gorm:"column:fats;default:0" json:"fats"`                   // Жиры (г)
	Carbohydrates  float64       `gorm:"column:carbohydrates;default:0" json:"carbohydrates"` // Углеводы (г)
	Calories       float64       `gorm:"column:calories;default:0" json:"calories"`           // Калории (ккал)
	Ingredients    []*Ingredient `gorm:"foreignKey:DishID" json:"ingredients"`
	Allergens      []*Allergen   `gorm:"foreignKey:DishID" json:"allergens"`
}

// Ingredient represents an ingredient in a dish
type Ingredient struct {
	gorm.Model
	DishID   uint    `gorm:"column:dish_id;not null" json:"dish_id"`
	Name     string  `gorm:"column:name;not null" json:"name"`
	Quantity float64 `gorm:"column:quantity;not null" json:"quantity"`
}

// Allergen represents an allergen in a dish
// Используется для информирования гостей об аллергенах в блюде
type Allergen struct {
	gorm.Model
	DishID      uint   `gorm:"column:dish_id;not null" json:"dish_id"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}
