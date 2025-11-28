package payload

// IngredientReq представляет ингредиент в запросе.
type IngredientReq struct {
	Name     string  `json:"name" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
}

// AllergenReq представляет аллерген в запросе.
type AllergenReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateDishReq запрос на создание нового блюда.
type CreateDishReq struct {
	RestaurantID   uint            `json:"restaurant_id" validate:"required"`
	MenuCategoryID uint            `json:"menu_category_id" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
	Allergens      []AllergenReq   `json:"allergens,omitempty" validate:"omitempty,dive"`
}

// UpdateDishReq запрос на обновление блюда.
type UpdateDishReq struct {
	RestaurantID   uint            `json:"restaurant_id" validate:"required"`
	MenuCategoryID uint            `json:"menu_category_id" validate:"required"`
	Name           string          `json:"name" validate:"required"`
	Price          float64         `json:"price" validate:"required"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	Ingredients    []IngredientReq `json:"ingredients" validate:"required,dive"`
	Allergens      []AllergenReq   `json:"allergens,omitempty" validate:"omitempty,dive"`
}
