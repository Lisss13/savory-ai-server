package payload

import (
	"time"
)

// IngredientResp represents an ingredient in a dish response
type IngredientResp struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

// OrganizationResp represents an organization in a dish response
type OrganizationResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// MenuCategoryResp represents a menu category in a dish response
type MenuCategoryResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// DishResp represents a dish in a response
type DishResp struct {
	ID           uint             `json:"id"`
	CreatedAt    time.Time        `json:"created_at"`
	Organization OrganizationResp `json:"organization"`
	MenuCategory MenuCategoryResp `json:"menu_category"`
	Name         string           `json:"name"`
	Price        float64          `json:"price"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	Ingredients  []IngredientResp `json:"ingredients"`
}

// DishesResp represents a list of dishes in a response
type DishesResp struct {
	Dishes []DishResp `json:"dishes"`
}

type DishCategoryResp struct {
	Category MenuCategoryResp `json:"category"`
	Dishes   []DishResp       `json:"dishes"`
}

type DishByCategoryResp struct {
	Dishes []DishCategoryResp `json:"dishes"`
}
