package payload

import "time"

// UserResp represents a user in a menu category response
type UserResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MenuCategoryResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

type MenuCategoriesResp struct {
	Categories []MenuCategoryResp `json:"categories"`
}
