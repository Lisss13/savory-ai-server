package payload

// UpdateUserStatusReq - запрос на изменение статуса пользователя
type UpdateUserStatusReq struct {
	IsActive bool `json:"isActive"`
}

// UpdateUserRoleReq - запрос на изменение роли пользователя
type UpdateUserRoleReq struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}

// CreateAdminUserReq - запрос на создание администратора
type CreateAdminUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateOrganizationStatusReq - запрос на изменение статуса организации
type UpdateOrganizationStatusReq struct {
	IsActive bool `json:"isActive"`
}

// ModerateContentReq - запрос на модерацию контента
type ModerateContentReq struct {
	Status string `json:"status" validate:"required,oneof=approved rejected"`
	Reason string `json:"reason,omitempty"`
}
