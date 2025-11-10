package payload

// IngredientReq represents an ingredient in a dish request
type IngredientReq struct {
	Name     string  `json:"name" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
}

// CreateDishReq represents a request to create a new dish
type CreateDishReq struct {
	MenuCategoryID uint            `json:"menu_category_id" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
}

// UpdateDishReq represents a request to update an existing dish
type UpdateDishReq struct {
	MenuCategoryID uint            `json:"menu_category_id" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
}
