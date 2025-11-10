package payload

import "time"

type UserResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Company   string    `json:"company"`
	Phone     string    `json:"phone"`
}
