package payload

import "time"

// WorkingHourResp represents a working hour in a restaurant response
type WorkingHourResp struct {
	ID        uint   `json:"id"`
	DayOfWeek int    `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// OrganizationResp represents a user in a restaurant response
type OrganizationResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// RestaurantResp represents a restaurant in a response
type RestaurantResp struct {
	ID           uint              `json:"id"`
	CreatedAt    time.Time         `json:"created_at"`
	Organization OrganizationResp  `json:"organization"`
	Name         string            `json:"name"`
	Address      string            `json:"address"`
	Phone        string            `json:"phone"`
	Website      string            `json:"website"`
	Description  string            `json:"description"`
	ImageURL     string            `json:"image_url"`
	WorkingHours []WorkingHourResp `json:"working_hours"`
}

// RestaurantsResp represents a list of restaurants in a response
type RestaurantsResp struct {
	Restaurants []RestaurantResp `json:"restaurants"`
}

// DeleteRestaurantResp represents a response after deleting a restaurant
type DeleteRestaurantResp struct {
	ID uint `json:"id"`
}
