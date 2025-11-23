package payload

import "time"

// CreateSubscriptionReq represents a request to create a new subscription
type CreateSubscriptionReq struct {
	OrganizationID uint      `json:"organizationId" validate:"required"`
	Period         int       `json:"period" validate:"required,min=1"` // Number of months
	StartDate      time.Time `json:"startDate" validate:"required"`
}

// UpdateSubscriptionReq represents a request to update an existing subscription
type UpdateSubscriptionReq struct {
	Period    int       `json:"period" validate:"required,min=1"` // Number of months
	StartDate time.Time `json:"startDate" validate:"required"`
	IsActive  bool      `json:"isActive"`
}

// ExtendSubscriptionReq represents a request to extend subscription period
type ExtendSubscriptionReq struct {
	Period int `json:"period" validate:"required,min=1"` // Additional months to add
}
