package payload

import "time"

// MenuCategoryResp ответ с данными категории меню.
type MenuCategoryResp struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	RestaurantID uint      `json:"restaurant_id"`
	SortOrder    int       `json:"sort_order"` // Порядок отображения категории
}

// MenuCategoriesResp ответ со списком категорий меню.
// Категории отсортированы по полю sort_order.
type MenuCategoriesResp struct {
	Categories []MenuCategoryResp `json:"categories"`
}
