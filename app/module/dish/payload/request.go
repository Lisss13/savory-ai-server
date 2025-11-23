package payload

// IngredientReq represents an ingredient in a dish request
type IngredientReq struct {
	Name     string  `json:"name" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
}

// AllergenReq represents an allergen in a dish request
// Используется для указания аллергенов при создании/обновлении блюда
type AllergenReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateDishReq represents a request to create a new dish
type CreateDishReq struct {
	MenuCategoryID uint            `json:"menuCategoryId" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
	Allergens      []AllergenReq   `json:"allergens,omitempty" validate:"omitempty,dive"`
}

// UpdateDishReq represents a request to update an existing dish
type UpdateDishReq struct {
	MenuCategoryID uint            `json:"menuCategoryId" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
	Allergens      []AllergenReq   `json:"allergens,omitempty" validate:"omitempty,dive"`
}
