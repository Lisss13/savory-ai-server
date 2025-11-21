package payload

type LoginResponse struct {
	Token        string                   `json:"token"`
	Type         string                   `json:"type"`
	ExpiresAt    int64                    `json:"expires_at"`
	User         RegisterResponse         `json:"user"`
	Organization UserOrganizationResponse `json:"organization"`
}

type UserOrganizationResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
	AdminID uint   `json:"admin_id"`
}

type RegisterResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}
