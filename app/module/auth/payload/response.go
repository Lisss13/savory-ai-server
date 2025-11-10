package payload

type LoginResponse struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	ExpiresAt int64  `json:"expires_at"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
