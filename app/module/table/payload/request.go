package payload

// CreateTableReq represents a request to create a new table
type CreateTableReq struct {
	Name         string `json:"name" validate:"required"`
	GuestCount   int    `json:"guestCount" validate:"required,min=1"`
	RestaurantID uint   `json:"restaurantId" validate:"required"`
}

// UpdateTableReq represents a request to update an existing table
type UpdateTableReq struct {
	Name         string `json:"name" validate:"required"`
	GuestCount   int    `json:"guestCount" validate:"required,min=1"`
	RestaurantID uint   `json:"restaurantId" validate:"required"`
}
