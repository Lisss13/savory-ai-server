package payload

// CreateOrganizationReq represents a request to create an organization
type CreateOrganizationReq struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	AdminID uint   `json:"admin_id" validate:"required"`
}

// UpdateOrganizationReq represents a request to update an organization
type UpdateOrganizationReq struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// AddUserToOrgReq represents a request to add a user to an organization
type AddUserToOrgReq struct {
	UserID uint `json:"user_id" validate:"required"`
}

// RemoveUserFromOrgReq represents a request to remove a user from an organization
type RemoveUserFromOrgReq struct {
	UserID uint `json:"user_id" validate:"required"`
}