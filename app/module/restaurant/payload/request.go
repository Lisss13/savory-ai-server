package payload

import "time"

// WorkingHourReq represents a working hour in a restaurant request
type WorkingHourReq struct {
	DayOfWeek int       `json:"day_of_week" validate:"required,min=0,max=6"`
	OpenTime  time.Time `json:"open_time" validate:"required"`
	CloseTime time.Time `json:"close_time" validate:"required"`
}

// CreateRestaurantReq represents a request to create a new restaurant
type CreateRestaurantReq struct {
	OrganizationID uint             `json:"organization_id" validate:"required"`
	Name           string           `json:"name" validate:"required"`
	Address        string           `json:"address" validate:"required"`
	Phone          string           `json:"phone" validate:"required"`
	Website        string           `json:"website"`
	Description    string           `json:"description"`
	ImageURL       string           `json:"image_url"`
	WorkingHours   []WorkingHourReq `json:"working_hours" validate:"required"`
}

// UpdateRestaurantReq represents a request to update an existing restaurant
type UpdateRestaurantReq struct {
	OrganizationID uint   `json:"organization_id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Address        string `json:"address" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Website        string `json:"website"`
	Description    string `json:"description"`
	ImageURL       string `json:"image_url"`
}
