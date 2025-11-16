package payload

import "time"

// RestaurantResp represents a restaurant in a table response
type RestaurantResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TableResp represents a table in a response
type TableResp struct {
	ID         uint           `json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	Restaurant RestaurantResp `json:"restaurant"`
	Name       string         `json:"name"`
	GuestCount int            `json:"guestCount"`
}

// TablesResp represents a list of tables in a response
type TablesResp struct {
	Tables []TableResp `json:"tables"`
}
