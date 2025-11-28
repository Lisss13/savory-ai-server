package payload

import (
	"time"
)

// IngredientResp представляет ингредиент в ответе.
type IngredientResp struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

// AllergenResp представляет аллерген в ответе.
type AllergenResp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// RestaurantResp представляет ресторан в ответе блюда.
type RestaurantResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// MenuCategoryResp представляет категорию меню в ответе блюда.
type MenuCategoryResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// DishResp представляет блюдо в ответе.
type DishResp struct {
	ID           uint             `json:"id"`
	CreatedAt    time.Time        `json:"created_at"`
	Restaurant   RestaurantResp   `json:"restaurant"`
	MenuCategory MenuCategoryResp `json:"menu_category"`
	Name         string           `json:"name"`
	Price        float64          `json:"price"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	Ingredients  []IngredientResp `json:"ingredients"`
	Allergens    []AllergenResp   `json:"allergens"`
}

// DishesResp представляет список блюд в ответе.
type DishesResp struct {
	Dishes []DishResp `json:"dishes"`
}

// DishCategoryResp представляет категорию с блюдами.
type DishCategoryResp struct {
	Category MenuCategoryResp `json:"category"`
	Dishes   []DishResp       `json:"dishes"`
}

// DishByCategoryResp представляет блюда, сгруппированные по категориям.
type DishByCategoryResp struct {
	Dishes []DishCategoryResp `json:"dishes"`
}
