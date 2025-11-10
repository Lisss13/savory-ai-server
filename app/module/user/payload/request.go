package payload

type UserUpdateReq struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email" validate:"email"`
	Phone   string `json:"phone"`
}

type UserCreateReq struct {
	Name     string `json:"name" validate:"required"`
	Company  string `json:"company" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=6"`
}
