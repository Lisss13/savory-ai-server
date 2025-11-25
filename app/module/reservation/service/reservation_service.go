// Package service содержит бизнес-логику для модуля бронирования столиков.
// Включает расчёт доступных слотов, проверку конфликтов и управление бронированиями.
package service

import (
	"errors"
	"fmt"
	"savory-ai-server/app/module/reservation/payload"
	"savory-ai-server/app/module/reservation/repository"
	restaurantRepo "savory-ai-server/app/module/restaurant/repository"
	tableRepo "savory-ai-server/app/module/table/repository"
	"savory-ai-server/app/storage"
	"sort"
	"time"
)

const (
	// SlotInterval определяет интервал между слотами бронирования в минутах.
	// Например, при значении 30 слоты будут: 12:00, 12:30, 13:00, и т.д.
	SlotInterval = 30
)

// reservationService реализует интерфейс ReservationService.
// Содержит зависимости от репозиториев бронирований, ресторанов и столиков.
type reservationService struct {
	reservationRepo repository.ReservationRepository
	restaurantRepo  restaurantRepo.RestaurantRepository
	tableRepo       tableRepo.TableRepository
}

// ReservationService определяет интерфейс бизнес-логики для бронирования.
// Разделён на методы для персонала и публичные методы для анонимных пользователей.
type ReservationService interface {
	// Методы для персонала ресторана
	GetAll() (*payload.ReservationsResp, error)                                                         // Получить все бронирования
	GetByID(id uint) (*payload.ReservationResp, error)                                                  // Получить бронь по ID
	GetByRestaurantID(restaurantID uint) (*payload.ReservationsResp, error)                             // Получить брони ресторана
	Update(id uint, req *payload.UpdateReservationReq) (*payload.ReservationResp, error)                // Обновить бронирование
	Cancel(id uint) (*payload.DeleteReservationResp, error)                                             // Отменить бронь (админ)
	Delete(id uint) (*payload.DeleteReservationResp, error)                                             // Удалить бронирование

	// Публичные методы для анонимных пользователей (через AI-чат)
	GetByPhone(phone string) (*payload.ReservationsResp, error)                                         // Получить брони по телефону
	GetAvailableSlots(restaurantID uint, date string, guestCount int) (*payload.AvailableSlotsResp, error) // Получить доступные слоты
	Create(req *payload.CreateReservationReq) (*payload.ReservationResp, error)                         // Создать бронирование
	CancelByPhone(id uint, phone string) (*payload.DeleteReservationResp, error)                        // Отменить бронь по телефону
}

// NewReservationService создаёт новый экземпляр сервиса бронирования.
// Требует репозитории для бронирований, ресторанов и столиков.
func NewReservationService(
	reservationRepo repository.ReservationRepository,
	restaurantRepo restaurantRepo.RestaurantRepository,
	tableRepo tableRepo.TableRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		restaurantRepo:  restaurantRepo,
		tableRepo:       tableRepo,
	}
}

// GetAll возвращает все бронирования из базы данных.
// Используется персоналом для просмотра всех бронирований.
func (s *reservationService) GetAll() (*payload.ReservationsResp, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var reservationResps []payload.ReservationResp
	for _, res := range reservations {
		reservationResps = append(reservationResps, mapReservationToResponse(&res))
	}

	return &payload.ReservationsResp{Reservations: reservationResps}, nil
}

// GetByID возвращает бронирование по его ID.
// Возвращает ошибку если бронирование не найдено.
func (s *reservationService) GetByID(id uint) (*payload.ReservationResp, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapReservationToResponse(reservation)
	return &resp, nil
}

// GetByPhone возвращает бронирования по номеру телефона клиента.
// Публичный метод для анонимных пользователей.
// Исключает отменённые бронирования.
func (s *reservationService) GetByPhone(phone string) (*payload.ReservationsResp, error) {
	reservations, err := s.reservationRepo.FindByPhone(phone)
	if err != nil {
		return nil, err
	}

	var reservationResps []payload.ReservationResp
	for _, res := range reservations {
		reservationResps = append(reservationResps, mapReservationToResponse(&res))
	}

	return &payload.ReservationsResp{Reservations: reservationResps}, nil
}

// GetByRestaurantID возвращает все бронирования для указанного ресторана.
// Сортировка по дате и времени (новые сначала).
func (s *reservationService) GetByRestaurantID(restaurantID uint) (*payload.ReservationsResp, error) {
	reservations, err := s.reservationRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var reservationResps []payload.ReservationResp
	for _, res := range reservations {
		reservationResps = append(reservationResps, mapReservationToResponse(&res))
	}

	return &payload.ReservationsResp{Reservations: reservationResps}, nil
}

// GetAvailableSlots возвращает доступные временные слоты для бронирования.
// Публичный метод для AI-чата.
//
// Алгоритм расчёта слотов:
//  1. Получает рабочие часы ресторана для указанного дня недели
//  2. Находит столики с подходящей вместимостью (>= guestCount)
//  3. Получает существующие бронирования на эту дату
//  4. Генерирует слоты с интервалом SlotInterval минут
//  5. Исключает слоты, которые конфликтуют с существующими бронями
//  6. Для текущего дня исключает прошедшие слоты
//
// Параметры:
//   - restaurantID: ID ресторана
//   - dateStr: дата в формате "YYYY-MM-DD"
//   - guestCount: количество гостей
//
// Возвращает список доступных слотов с информацией о столиках.
func (s *reservationService) GetAvailableSlots(restaurantID uint, dateStr string, guestCount int) (*payload.AvailableSlotsResp, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Get restaurant with working hours
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Get day of week (0 = Sunday)
	dayOfWeek := int(date.Weekday())

	// Find working hours for this day
	var workingHour *storage.WorkingHour
	for _, wh := range restaurant.WorkingHours {
		if wh.DayOfWeek == dayOfWeek {
			workingHour = wh
			break
		}
	}

	if workingHour == nil {
		return &payload.AvailableSlotsResp{
			Date:  dateStr,
			Slots: []payload.TimeSlot{},
		}, nil // Restaurant is closed on this day
	}

	// Get tables that can accommodate the guest count
	tables, err := s.tableRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var suitableTables []*storage.Table
	for _, table := range tables {
		if table.GuestCount >= guestCount {
			suitableTables = append(suitableTables, table)
		}
	}

	if len(suitableTables) == 0 {
		return &payload.AvailableSlotsResp{
			Date:  dateStr,
			Slots: []payload.TimeSlot{},
		}, nil // No tables can accommodate the guest count
	}

	// Get existing reservations for this date
	existingReservations, err := s.reservationRepo.FindByRestaurantIDAndDate(restaurantID, date)
	if err != nil {
		return nil, err
	}

	// Build map of table -> reserved time slots
	tableReservations := make(map[uint][]timeRange)
	for _, res := range existingReservations {
		tableReservations[res.TableID] = append(tableReservations[res.TableID], timeRange{
			start: res.StartTime,
			end:   res.EndTime,
		})
	}

	// Generate available slots
	duration := restaurant.ReservationDuration
	if duration == 0 {
		duration = 90 // default 90 minutes
	}

	openTime, _ := time.Parse("15:04", workingHour.OpenTime)
	closeTime, _ := time.Parse("15:04", workingHour.CloseTime)

	var slots []payload.TimeSlot

	// Generate slots for each suitable table
	for _, table := range suitableTables {
		currentTime := openTime
		for currentTime.Add(time.Duration(duration) * time.Minute).Before(closeTime) ||
			currentTime.Add(time.Duration(duration)*time.Minute).Equal(closeTime) {

			startTimeStr := currentTime.Format("15:04")
			endTime := currentTime.Add(time.Duration(duration) * time.Minute)
			endTimeStr := endTime.Format("15:04")

			// Check if this slot conflicts with existing reservations
			if !hasConflict(tableReservations[table.ID], startTimeStr, endTimeStr) {
				// Check if slot is not in the past (for today)
				if date.After(time.Now()) || (date.Day() == time.Now().Day() && date.Month() == time.Now().Month() && date.Year() == time.Now().Year()) {
					now := time.Now()
					slotDateTime := time.Date(date.Year(), date.Month(), date.Day(),
						currentTime.Hour(), currentTime.Minute(), 0, 0, time.Local)
					if slotDateTime.After(now) {
						slots = append(slots, payload.TimeSlot{
							StartTime: startTimeStr,
							EndTime:   endTimeStr,
							TableID:   table.ID,
							TableName: table.Name,
							Capacity:  table.GuestCount,
						})
					}
				} else {
					slots = append(slots, payload.TimeSlot{
						StartTime: startTimeStr,
						EndTime:   endTimeStr,
						TableID:   table.ID,
						TableName: table.Name,
						Capacity:  table.GuestCount,
					})
				}
			}

			currentTime = currentTime.Add(time.Duration(SlotInterval) * time.Minute)
		}
	}

	// Sort slots by start time
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartTime < slots[j].StartTime
	})

	return &payload.AvailableSlotsResp{
		Date:  dateStr,
		Slots: slots,
	}, nil
}

// Create создаёт новое бронирование столика.
// Публичный метод для AI-чата.
//
// Процесс создания:
//  1. Валидирует формат даты и времени
//  2. Проверяет существование ресторана и столика
//  3. Проверяет принадлежность столика ресторану
//  4. Рассчитывает время окончания на основе ReservationDuration ресторана
//  5. Проверяет наличие конфликтов с существующими бронями
//  6. Создаёт бронирование со статусом "confirmed"
//
// Параметры req:
//   - RestaurantID, TableID: обязательные ID ресторана и столика
//   - CustomerName, CustomerPhone: данные клиента (обязательные)
//   - GuestCount: количество гостей
//   - ReservationDate: дата в формате "YYYY-MM-DD"
//   - StartTime: время начала в формате "HH:MM"
//
// Возвращает созданное бронирование или ошибку при конфликте.
func (s *reservationService) Create(req *payload.CreateReservationReq) (*payload.ReservationResp, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.ReservationDate)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Validate start time format
	_, err = time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, errors.New("invalid time format, use HH:MM")
	}

	// Get restaurant for duration
	restaurant, err := s.restaurantRepo.FindByID(req.RestaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Validate table exists and belongs to restaurant
	table, err := s.tableRepo.FindByID(req.TableID)
	if err != nil {
		return nil, errors.New("table not found")
	}
	if table.RestaurantID != req.RestaurantID {
		return nil, errors.New("table does not belong to this restaurant")
	}

	// Calculate end time
	duration := restaurant.ReservationDuration
	if duration == 0 {
		duration = 90
	}
	startTime, _ := time.Parse("15:04", req.StartTime)
	endTime := startTime.Add(time.Duration(duration) * time.Minute)
	endTimeStr := endTime.Format("15:04")

	// Check for conflicts
	existingReservations, err := s.reservationRepo.FindByTableIDAndDate(req.TableID, date)
	if err != nil {
		return nil, err
	}

	var tableReservations []timeRange
	for _, res := range existingReservations {
		tableReservations = append(tableReservations, timeRange{
			start: res.StartTime,
			end:   res.EndTime,
		})
	}

	if hasConflict(tableReservations, req.StartTime, endTimeStr) {
		return nil, errors.New("this time slot is already booked")
	}

	// Create reservation
	reservation := &storage.Reservation{
		RestaurantID:    req.RestaurantID,
		TableID:         req.TableID,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerEmail:   req.CustomerEmail,
		GuestCount:      req.GuestCount,
		ReservationDate: date,
		StartTime:       req.StartTime,
		EndTime:         endTimeStr,
		Status:          storage.ReservationStatusConfirmed,
		Notes:           req.Notes,
		ChatSessionID:   req.ChatSessionID,
	}

	createdReservation, err := s.reservationRepo.Create(reservation)
	if err != nil {
		return nil, err
	}

	resp := mapReservationToResponse(createdReservation)
	return &resp, nil
}

// Update обновляет существующее бронирование.
// Позволяет изменить данные клиента, количество гостей, дату и время.
// При изменении даты/времени проверяет конфликты с другими бронями.
//
// Поддерживает частичное обновление - можно передать только изменяемые поля.
func (s *reservationService) Update(id uint, req *payload.UpdateReservationReq) (*payload.ReservationResp, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("reservation not found")
	}

	if req.CustomerName != "" {
		reservation.CustomerName = req.CustomerName
	}
	if req.CustomerPhone != "" {
		reservation.CustomerPhone = req.CustomerPhone
	}
	if req.CustomerEmail != "" {
		reservation.CustomerEmail = req.CustomerEmail
	}
	if req.GuestCount > 0 {
		reservation.GuestCount = req.GuestCount
	}
	if req.Notes != "" {
		reservation.Notes = req.Notes
	}

	// Handle date/time update
	if req.ReservationDate != "" || req.StartTime != "" {
		newDate := reservation.ReservationDate
		newStartTime := reservation.StartTime

		if req.ReservationDate != "" {
			parsedDate, err := time.Parse("2006-01-02", req.ReservationDate)
			if err != nil {
				return nil, errors.New("invalid date format")
			}
			newDate = parsedDate
		}

		if req.StartTime != "" {
			_, err := time.Parse("15:04", req.StartTime)
			if err != nil {
				return nil, errors.New("invalid time format")
			}
			newStartTime = req.StartTime
		}

		// Recalculate end time
		restaurant, err := s.restaurantRepo.FindByID(reservation.RestaurantID)
		if err != nil {
			return nil, err
		}

		duration := restaurant.ReservationDuration
		if duration == 0 {
			duration = 90
		}
		startTime, _ := time.Parse("15:04", newStartTime)
		endTime := startTime.Add(time.Duration(duration) * time.Minute)

		// Check for conflicts (excluding current reservation)
		existingReservations, err := s.reservationRepo.FindByTableIDAndDate(reservation.TableID, newDate)
		if err != nil {
			return nil, err
		}

		var tableReservations []timeRange
		for _, res := range existingReservations {
			if res.ID != id { // Exclude current reservation
				tableReservations = append(tableReservations, timeRange{
					start: res.StartTime,
					end:   res.EndTime,
				})
			}
		}

		if hasConflict(tableReservations, newStartTime, endTime.Format("15:04")) {
			return nil, errors.New("this time slot is already booked")
		}

		reservation.ReservationDate = newDate
		reservation.StartTime = newStartTime
		reservation.EndTime = endTime.Format("15:04")
	}

	updatedReservation, err := s.reservationRepo.Update(reservation)
	if err != nil {
		return nil, err
	}

	resp := mapReservationToResponse(updatedReservation)
	return &resp, nil
}

// Cancel отменяет бронирование (административная отмена).
// Меняет статус бронирования на "cancelled".
// Запись не удаляется из базы данных.
func (s *reservationService) Cancel(id uint) (*payload.DeleteReservationResp, error) {
	err := s.reservationRepo.UpdateStatus(id, storage.ReservationStatusCancelled)
	if err != nil {
		return nil, err
	}

	return &payload.DeleteReservationResp{
		ID:      id,
		Message: "Reservation cancelled successfully",
	}, nil
}

// CancelByPhone отменяет бронирование с верификацией по номеру телефона.
// Публичный метод для анонимных пользователей (через AI-чат).
//
// Проверки:
//   - Бронирование существует и принадлежит указанному телефону
//   - Бронирование ещё не отменено
//
// Возвращает ошибку если телефон не соответствует бронированию.
func (s *reservationService) CancelByPhone(id uint, phone string) (*payload.DeleteReservationResp, error) {
	// Verify the reservation belongs to this phone number
	reservation, err := s.reservationRepo.FindByIDAndPhone(id, phone)
	if err != nil {
		return nil, errors.New("reservation not found or phone number doesn't match")
	}

	if reservation.Status == storage.ReservationStatusCancelled {
		return nil, errors.New("reservation is already cancelled")
	}

	err = s.reservationRepo.UpdateStatus(id, storage.ReservationStatusCancelled)
	if err != nil {
		return nil, err
	}

	return &payload.DeleteReservationResp{
		ID:      id,
		Message: "Reservation cancelled successfully",
	}, nil
}

// Delete полностью удаляет бронирование из базы данных.
// В отличие от Cancel, запись удаляется безвозвратно.
func (s *reservationService) Delete(id uint) (*payload.DeleteReservationResp, error) {
	err := s.reservationRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	return &payload.DeleteReservationResp{
		ID:      id,
		Message: "Reservation deleted successfully",
	}, nil
}

// =====================================================
// Вспомогательные типы и функции
// =====================================================

// timeRange представляет временной диапазон для проверки конфликтов.
type timeRange struct {
	start string // Время начала в формате "HH:MM"
	end   string // Время окончания в формате "HH:MM"
}

// hasConflict проверяет, пересекается ли новый слот с существующими бронированиями.
// Использует простую проверку пересечения интервалов:
// конфликт есть, если НЕ выполняется условие (newEnd <= existingStart ИЛИ newStart >= existingEnd)
func hasConflict(reservations []timeRange, newStart, newEnd string) bool {
	for _, r := range reservations {
		// Check if times overlap
		if !(newEnd <= r.start || newStart >= r.end) {
			return true
		}
	}
	return false
}

// mapReservationToResponse преобразует модель Reservation в DTO для ответа API.
// Включает данные связанных сущностей (ресторан, столик).
func mapReservationToResponse(res *storage.Reservation) payload.ReservationResp {
	return payload.ReservationResp{
		ID:              res.ID,
		RestaurantID:    res.RestaurantID,
		RestaurantName:  res.Restaurant.Name,
		TableID:         res.TableID,
		TableName:       res.Table.Name,
		CustomerName:    res.CustomerName,
		CustomerPhone:   res.CustomerPhone,
		CustomerEmail:   res.CustomerEmail,
		GuestCount:      res.GuestCount,
		ReservationDate: res.ReservationDate.Format("2006-01-02"),
		StartTime:       res.StartTime,
		EndTime:         res.EndTime,
		Status:          string(res.Status),
		Notes:           res.Notes,
		CreatedAt:       res.CreatedAt,
	}
}

// GetAvailableSlotsForAI возвращает доступные слоты в текстовом формате для AI-ответа.
// Используется Anthropic сервисом для формирования человекочитаемого ответа.
func (s *reservationService) GetAvailableSlotsForAI(restaurantID uint, dateStr string, guestCount int) (string, error) {
	slots, err := s.GetAvailableSlots(restaurantID, dateStr, guestCount)
	if err != nil {
		return "", err
	}

	if len(slots.Slots) == 0 {
		return fmt.Sprintf("No available slots on %s for %d guests.", dateStr, guestCount), nil
	}

	result := fmt.Sprintf("Available slots on %s for %d guests:\n", dateStr, guestCount)
	for _, slot := range slots.Slots {
		result += fmt.Sprintf("- %s - %s (Table: %s, capacity: %d)\n",
			slot.StartTime, slot.EndTime, slot.TableName, slot.Capacity)
	}

	return result, nil
}
