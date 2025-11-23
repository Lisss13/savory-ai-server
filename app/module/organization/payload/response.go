package payload

import (
	"time"
)

// UserInOrgResp represents a user in an organization response
type UserInOrgResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// OrganizationResp represents an organization response
type OrganizationResp struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Name      string          `json:"name"`
	Phone     string          `json:"phone"`
	Admin     UserInOrgResp   `json:"admin"`
	Users     []UserInOrgResp `json:"users,omitempty"`
	Languages []LanguageResp  `json:"languages"`
}

// OrganizationsResp represents a list of organizations
type OrganizationsResp struct {
	Organizations []OrganizationResp `json:"organizations"`
}
