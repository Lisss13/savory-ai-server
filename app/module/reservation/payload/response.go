package payload

import "time"

// ReservationResp represents a reservation in a response
type ReservationResp struct {
	ID              uint      `json:"id"`
	RestaurantID    uint      `json:"restaurant_id"`
	RestaurantName  string    `json:"restaurant_name"`
	TableID         uint      `json:"table_id"`
	TableName       string    `json:"table_name"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerEmail   string    `json:"customer_email,omitempty"`
	GuestCount      int       `json:"guest_count"`
	ReservationDate string    `json:"reservation_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// ReservationsResp represents a list of reservations
type ReservationsResp struct {
	Reservations []ReservationResp `json:"reservations"`
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	TableID   uint   `json:"table_id"`
	TableName string `json:"table_name"`
	Capacity  int    `json:"capacity"`
}

// AvailableSlotsResp represents available time slots for a given date
type AvailableSlotsResp struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

// DeleteReservationResp represents a response after deleting/cancelling a reservation
type DeleteReservationResp struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
