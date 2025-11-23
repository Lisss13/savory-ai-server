package payload

import (
	"time"
)

// LanguageResp represents a language response
type LanguageResp struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// LanguagesResp represents a list of languages
type LanguagesResp struct {
	Languages []LanguageResp `json:"languages"`
}

// OrganizationLanguagesResp represents the languages for an organization
type OrganizationLanguagesResp struct {
	OrganizationID uint          `json:"organizationId"`
	Languages      []LanguageResp `json:"languages"`
}