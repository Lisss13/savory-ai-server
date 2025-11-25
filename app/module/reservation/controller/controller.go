package controller

// Controller агрегирует все контроллеры модуля бронирования.
// Используется для группировки контроллеров и передачи в роутер.
type Controller struct {
	Reservation ReservationController // Контроллер для работы с бронированиями
}

// NewControllers создаёт агрегатор контроллеров.
// Используется FX для dependency injection.
func NewControllers(reservation ReservationController) *Controller {
	return &Controller{Reservation: reservation}
}
