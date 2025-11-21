package payload

// StartTableSessionReq Создание чата для столика в ресторане
type StartTableSessionReq struct {
	TableID      uint `json:"tableId" validate:"required"`
	RestaurantID uint `json:"restaurantId" validate:"required"`
}

type CloseTableSessionReq struct {
	SessionID uint `json:"sessionId" validate:"required"`
}

type SendTableMessageReq struct {
	SessionID uint   `json:"sessionId" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// StartRestaurantSessionReq Создание чата для ресторана
type StartRestaurantSessionReq struct {
	RestaurantID uint `json:"restaurantId" validate:"required"`
}
