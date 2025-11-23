package payload

import "time"

// OrganizationResp represents an organization in a subscription response
type OrganizationResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// SubscriptionResp represents a subscription in a response
type SubscriptionResp struct {
	ID           uint             `json:"id"`
	CreatedAt    time.Time        `json:"createdAt"`
	Organization OrganizationResp `json:"organization"`
	Period       int              `json:"period"`    // Number of months
	StartDate    time.Time        `json:"startDate"` // Subscription start date
	EndDate      time.Time        `json:"endDate"`   // Subscription end date
	IsActive     bool             `json:"isActive"`  // Is subscription active
	DaysLeft     int              `json:"daysLeft"`  // Days left until expiration
}

// SubscriptionsResp represents a list of subscriptions in a response
type SubscriptionsResp struct {
	Subscriptions []SubscriptionResp `json:"subscriptions"`
}
