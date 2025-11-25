package payload

// CreateReservationReq represents a request to create a new reservation
type CreateReservationReq struct {
	RestaurantID    uint   `json:"restaurant_id" validate:"required"`
	TableID         uint   `json:"table_id" validate:"required"`
	CustomerName    string `json:"customer_name" validate:"required"`
	CustomerPhone   string `json:"customer_phone" validate:"required"`
	CustomerEmail   string `json:"customer_email" validate:"omitempty,email"`
	GuestCount      int    `json:"guest_count" validate:"required,min=1"`
	ReservationDate string `json:"reservation_date" validate:"required"` // Format: "2006-01-02"
	StartTime       string `json:"start_time" validate:"required"`       // Format: "15:04"
	Notes           string `json:"notes"`
	ChatSessionID   *uint  `json:"chat_session_id"`
}

// UpdateReservationReq represents a request to update an existing reservation
type UpdateReservationReq struct {
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerEmail   string `json:"customer_email" validate:"omitempty,email"`
	GuestCount      int    `json:"guest_count" validate:"omitempty,min=1"`
	ReservationDate string `json:"reservation_date"` // Format: "2006-01-02"
	StartTime       string `json:"start_time"`       // Format: "15:04"
	Notes           string `json:"notes"`
}

// CancelReservationReq represents a request to cancel a reservation
type CancelReservationReq struct {
	Reason string `json:"reason"`
}

// GetAvailableSlotsReq represents a request to get available time slots
type GetAvailableSlotsReq struct {
	RestaurantID uint   `json:"restaurant_id" validate:"required"`
	Date         string `json:"date" validate:"required"` // Format: "2006-01-02"
	GuestCount   int    `json:"guest_count" validate:"required,min=1"`
}
