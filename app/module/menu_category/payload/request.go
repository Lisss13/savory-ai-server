package payload

// CreateMenuCategoryReq запрос на создание категории меню.
type CreateMenuCategoryReq struct {
	Name         string `json:"name" validate:"required"`
	RestaurantID uint   `json:"restaurant_id" validate:"required"`
}
