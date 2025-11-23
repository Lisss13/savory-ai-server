package payload

// CreateLanguageReq represents a request to create a language
type CreateLanguageReq struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateLanguageReq represents a request to update a language
type UpdateLanguageReq struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddLanguageToOrgReq represents a request to add a language to an organization
type AddLanguageToOrgReq struct {
	LanguageID uint `json:"languageId" validate:"required"`
}

// RemoveLanguageFromOrgReq represents a request to remove a language from an organization
type RemoveLanguageFromOrgReq struct {
	LanguageID uint `json:"languageId" validate:"required"`
}