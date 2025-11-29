package payload

// CreateMenuCategoryReq запрос на создание категории меню.
type CreateMenuCategoryReq struct {
	Name         string `json:"name" validate:"required"`
	RestaurantID uint   `json:"restaurant_id" validate:"required"`
	SortOrder    *int   `json:"sort_order"` // Опционально, если не указан - будет установлен автоматически
}

// UpdateMenuCategoryReq запрос на обновление категории меню.
type UpdateMenuCategoryReq struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
}

// CategorySortOrderItem элемент для массового обновления порядка категорий.
type CategorySortOrderItem struct {
	ID        uint `json:"id" validate:"required"`
	SortOrder int  `json:"sort_order" validate:"min=0"`
}

// UpdateCategoriesSortOrderReq запрос на массовое обновление порядка категорий.
// Позволяет пользователю изменить порядок отображения нескольких категорий за один запрос.
type UpdateCategoriesSortOrderReq struct {
	Categories []CategorySortOrderItem `json:"categories" validate:"required,min=1,dive"`
}
