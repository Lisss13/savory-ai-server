// Package controller содержит HTTP-обработчики для модуля бронирования столиков.
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/reservation/payload"
	"savory-ai-server/app/module/reservation/service"
	"savory-ai-server/utils/response"
)

// reservationController реализует интерфейс ReservationController.
// Содержит зависимость от сервиса бронирования.
type reservationController struct {
	reservationService service.ReservationService
}

// ReservationController определяет интерфейс HTTP-обработчиков для бронирования.
// Включает методы для публичного доступа (анонимные пользователи) и защищённые (персонал).
type ReservationController interface {
	// Защищённые методы (требуют авторизации)
	GetAll(c *fiber.Ctx) error            // Получить все бронирования
	GetByID(c *fiber.Ctx) error           // Получить бронь по ID
	GetByRestaurantID(c *fiber.Ctx) error // Получить брони ресторана
	Update(c *fiber.Ctx) error            // Обновить бронирование
	Cancel(c *fiber.Ctx) error            // Отменить бронь (админ)
	Delete(c *fiber.Ctx) error            // Удалить бронирование

	// Публичные методы (для анонимных пользователей через AI-чат)
	GetByPhone(c *fiber.Ctx) error        // Получить брони по телефону
	GetAvailableSlots(c *fiber.Ctx) error // Получить доступные слоты
	Create(c *fiber.Ctx) error            // Создать бронирование
	CancelByPhone(c *fiber.Ctx) error     // Отменить бронь по телефону
}

// NewReservationController создаёт новый экземпляр контроллера бронирования.
func NewReservationController(reservationService service.ReservationService) ReservationController {
	return &reservationController{reservationService: reservationService}
}

// GetAll возвращает список всех бронирований.
// Защищённый эндпоинт для персонала ресторана.
//
// Метод: GET /reservations/
// Требует: Авторизация (JWT)
// Ответ: Список всех бронирований
func (ctrl *reservationController) GetAll(c *fiber.Ctx) error {
	data, err := ctrl.reservationService.GetAll()
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByID возвращает бронирование по его ID.
// Защищённый эндпоинт для персонала ресторана.
//
// Метод: GET /reservations/:id
// Параметры: id - ID бронирования
// Требует: Авторизация (JWT)
// Ответ: Данные бронирования или 404 если не найдено
func (ctrl *reservationController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.GetByID(uint(id))
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusNotFound,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByPhone возвращает бронирования по номеру телефона клиента.
// Публичный эндпоинт для анонимных пользователей (через AI-чат).
// Телефон используется как идентификатор клиента.
//
// Метод: GET /reservations/my?phone=+79161234567
// Query параметры: phone - номер телефона клиента (обязательный)
// Требует: Не требует авторизации
// Ответ: Список бронирований клиента (кроме отменённых)
func (ctrl *reservationController) GetByPhone(c *fiber.Ctx) error {
	phone := c.Query("phone")
	if phone == "" {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Phone number is required"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.GetByPhone(phone)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByRestaurantID возвращает все бронирования для указанного ресторана.
// Защищённый эндпоинт для персонала ресторана.
//
// Метод: GET /reservations/restaurant/:restaurant_id
// Параметры: restaurant_id - ID ресторана
// Требует: Авторизация (JWT)
// Ответ: Список бронирований ресторана
func (ctrl *reservationController) GetByRestaurantID(c *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(c.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.GetByRestaurantID(uint(restaurantID))
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetAvailableSlots возвращает доступные временные слоты для бронирования.
// Публичный эндпоинт для анонимных пользователей (через AI-чат).
// Учитывает рабочие часы ресторана, вместимость столов и существующие брони.
//
// Метод: GET /reservations/available/:restaurant_id?date=2025-11-28&guest_count=2
// Параметры: restaurant_id - ID ресторана
// Query параметры:
//   - date (обязательный) - дата в формате YYYY-MM-DD
//   - guest_count (опционально) - количество гостей, по умолчанию 2
//
// Требует: Не требует авторизации
// Ответ: Список доступных слотов с информацией о столиках
func (ctrl *reservationController) GetAvailableSlots(c *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(c.Params("restaurant_id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid restaurant ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	date := c.Query("date")
	if date == "" {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Date is required (format: YYYY-MM-DD)"},
			Code:     fiber.StatusBadRequest,
		})
	}

	guestCountStr := c.Query("guest_count", "2")
	guestCount, err := strconv.Atoi(guestCountStr)
	if err != nil || guestCount < 1 {
		guestCount = 2
	}

	data, err := ctrl.reservationService.GetAvailableSlots(uint(restaurantID), date, guestCount)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// Create создаёт новое бронирование столика.
// Публичный эндпоинт для анонимных пользователей (через AI-чат).
// Автоматически подтверждает бронь если слот свободен.
//
// Метод: POST /reservations/
// Тело запроса: CreateReservationReq (restaurant_id, table_id, customer_name, customer_phone, guest_count, reservation_date, start_time)
// Требует: Не требует авторизации
// Ответ: Созданное бронирование со статусом "confirmed"
func (ctrl *reservationController) Create(c *fiber.Ctx) error {
	req := new(payload.CreateReservationReq)
	if err := c.BodyParser(req); err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.Create(req)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"Reservation created successfully"},
		Code:     fiber.StatusCreated,
	})
}

// Update обновляет существующее бронирование.
// Защищённый эндпоинт для персонала ресторана.
// Позволяет изменить данные клиента, количество гостей, дату/время.
// При изменении даты/времени проверяет наличие конфликтов.
//
// Метод: PATCH /reservations/:id
// Параметры: id - ID бронирования
// Тело запроса: UpdateReservationReq (все поля опциональны)
// Требует: Авторизация (JWT)
// Ответ: Обновлённое бронирование
func (ctrl *reservationController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.UpdateReservationReq)
	if err := c.BodyParser(req); err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.Update(uint(id), req)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"Reservation updated successfully"},
		Code:     fiber.StatusOK,
	})
}

// Cancel отменяет бронирование (административная отмена).
// Защищённый эндпоинт для персонала ресторана.
// Меняет статус бронирования на "cancelled".
//
// Метод: POST /reservations/:id/cancel
// Параметры: id - ID бронирования
// Требует: Авторизация (JWT)
// Ответ: Подтверждение отмены
func (ctrl *reservationController) Cancel(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.Cancel(uint(id))
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"Reservation cancelled successfully"},
		Code:     fiber.StatusOK,
	})
}

// CancelByPhone отменяет бронирование с верификацией по номеру телефона.
// Публичный эндпоинт для анонимных пользователей (через AI-чат).
// Проверяет, что номер телефона соответствует бронированию.
//
// Метод: POST /reservations/:id/cancel/public?phone=+79161234567
// Параметры: id - ID бронирования
// Query параметры: phone - номер телефона для верификации (обязательный)
// Требует: Не требует авторизации
// Ответ: Подтверждение отмены или ошибка если телефон не совпадает
func (ctrl *reservationController) CancelByPhone(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.CancelReservationReq)
	if err := c.BodyParser(req); err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	phone := c.Query("phone")
	if phone == "" {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Phone number is required"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.CancelByPhone(uint(id), phone)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"Reservation cancelled successfully"},
		Code:     fiber.StatusOK,
	})
}

// Delete полностью удаляет бронирование из базы данных.
// Защищённый эндпоинт для персонала ресторана.
// В отличие от Cancel, полностью удаляет запись.
//
// Метод: DELETE /reservations/:id
// Параметры: id - ID бронирования
// Требует: Авторизация (JWT)
// Ответ: Подтверждение удаления
func (ctrl *reservationController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{"Invalid ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	data, err := ctrl.reservationService.Delete(uint(id))
	if err != nil {
		return response.Resp(c, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(c, response.Response{
		Data:     data,
		Messages: response.Messages{"Reservation deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
