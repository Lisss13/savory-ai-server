package payload

// CreateTableReq represents a request to create a new table
type CreateTableReq struct {
	Name       string `json:"name" validate:"required"`
	GuestCount int    `json:"guest_count" validate:"required,min=1"`
}

// UpdateTableReq represents a request to update an existing table
type UpdateTableReq struct {
	Name       string `json:"name" validate:"required"`
	GuestCount int    `json:"guest_count" validate:"required,min=1"`
}
