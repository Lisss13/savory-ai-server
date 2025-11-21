package payload

type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"12345678" validate:"required,min=8,max=255"`
}

type RegisterRequest struct {
	CompanyName string `json:"company" example:"Attic"`
	Name        string `json:"name" example:"John Doe"`
	Email       string `json:"email" example:"john.doe@gmail.com" validate:"required,email"`
	Phone       string `json:"phone" example:"+1234567890"`
	Password    string `json:"password" example:"12345678" validate:"required,min=8,max=255"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" example:"12345678" validate:"required,min=8,max=255"`
	NewPassword string `json:"newPassword" example:"87654321" validate:"required,min=8,max=255"`
}
