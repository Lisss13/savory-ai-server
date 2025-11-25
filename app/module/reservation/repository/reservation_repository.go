// Package repository содержит слой доступа к данным для модуля бронирования.
// Использует GORM для работы с PostgreSQL.
package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
	"time"
)

// reservationRepository реализует интерфейс ReservationRepository.
type reservationRepository struct {
	DB *database.Database
}

// ReservationRepository определяет интерфейс для работы с бронированиями в БД.
// Включает методы CRUD и специализированные запросы для бизнес-логики.
type ReservationRepository interface {
	// Методы поиска
	FindAll() ([]storage.Reservation, error)                                        // Все бронирования
	FindByID(id uint) (*storage.Reservation, error)                                 // По ID
	FindByIDAndPhone(id uint, phone string) (*storage.Reservation, error)           // По ID и телефону (верификация)
	FindByPhone(phone string) ([]storage.Reservation, error)                        // По телефону клиента
	FindByRestaurantID(restaurantID uint) ([]storage.Reservation, error)            // По ресторану
	FindByRestaurantIDAndDate(restaurantID uint, date time.Time) ([]storage.Reservation, error) // По ресторану и дате
	FindByTableIDAndDate(tableID uint, date time.Time) ([]storage.Reservation, error)           // По столику и дате

	// Методы изменения
	Create(reservation *storage.Reservation) (*storage.Reservation, error)          // Создать
	Update(reservation *storage.Reservation) (*storage.Reservation, error)          // Обновить
	Delete(id uint) error                                                            // Удалить
	UpdateStatus(id uint, status storage.ReservationStatus) error                    // Обновить статус
}

// NewReservationRepository создаёт новый экземпляр репозитория.
func NewReservationRepository(db *database.Database) ReservationRepository {
	return &reservationRepository{DB: db}
}

// FindAll возвращает все бронирования с предзагрузкой связей (Restaurant, Table).
func (r *reservationRepository) FindAll() ([]storage.Reservation, error) {
	var reservations []storage.Reservation
	err := r.DB.DB.Preload("Restaurant").Preload("Table").Find(&reservations).Error
	return reservations, err
}

// FindByID находит бронирование по ID с предзагрузкой связей.
func (r *reservationRepository) FindByID(id uint) (*storage.Reservation, error) {
	var reservation storage.Reservation
	err := r.DB.DB.Preload("Restaurant").Preload("Table").First(&reservation, id).Error
	return &reservation, err
}

// FindByIDAndPhone находит бронирование по ID и номеру телефона.
// Используется для верификации владельца при отмене через публичный API.
func (r *reservationRepository) FindByIDAndPhone(id uint, phone string) (*storage.Reservation, error) {
	var reservation storage.Reservation
	err := r.DB.DB.Preload("Restaurant").Preload("Table").
		Where("id = ? AND customer_phone = ?", id, phone).
		First(&reservation).Error
	return &reservation, err
}

// FindByPhone находит все бронирования клиента по номеру телефона.
// Исключает отменённые бронирования.
// Сортировка: новые бронирования сначала.
func (r *reservationRepository) FindByPhone(phone string) ([]storage.Reservation, error) {
	var reservations []storage.Reservation
	err := r.DB.DB.Preload("Restaurant").Preload("Table").
		Where("customer_phone = ? AND status NOT IN ?", phone, []storage.ReservationStatus{storage.ReservationStatusCancelled}).
		Order("reservation_date DESC, start_time DESC").
		Find(&reservations).Error
	return reservations, err
}

// FindByRestaurantID находит все бронирования для указанного ресторана.
// Сортировка: новые бронирования сначала.
func (r *reservationRepository) FindByRestaurantID(restaurantID uint) ([]storage.Reservation, error) {
	var reservations []storage.Reservation
	err := r.DB.DB.Preload("Restaurant").Preload("Table").
		Where("restaurant_id = ?", restaurantID).
		Order("reservation_date DESC, start_time DESC").
		Find(&reservations).Error
	return reservations, err
}

// FindByRestaurantIDAndDate находит бронирования ресторана на указанную дату.
// Используется для расчёта доступных слотов.
// Исключает отменённые бронирования.
// Сортировка по времени начала.
func (r *reservationRepository) FindByRestaurantIDAndDate(restaurantID uint, date time.Time) ([]storage.Reservation, error) {
	var reservations []storage.Reservation
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.DB.DB.Preload("Restaurant").Preload("Table").
		Where("restaurant_id = ? AND reservation_date >= ? AND reservation_date < ? AND status NOT IN ?",
			restaurantID, startOfDay, endOfDay, []storage.ReservationStatus{storage.ReservationStatusCancelled}).
		Order("start_time ASC").
		Find(&reservations).Error
	return reservations, err
}

// FindByTableIDAndDate находит бронирования столика на указанную дату.
// Используется для проверки конфликтов при создании/обновлении брони.
// Исключает отменённые бронирования.
func (r *reservationRepository) FindByTableIDAndDate(tableID uint, date time.Time) ([]storage.Reservation, error) {
	var reservations []storage.Reservation
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.DB.DB.Preload("Restaurant").Preload("Table").
		Where("table_id = ? AND reservation_date >= ? AND reservation_date < ? AND status NOT IN ?",
			tableID, startOfDay, endOfDay, []storage.ReservationStatus{storage.ReservationStatusCancelled}).
		Order("start_time ASC").
		Find(&reservations).Error
	return reservations, err
}

// Create создаёт новое бронирование и возвращает его с предзагруженными связями.
func (r *reservationRepository) Create(reservation *storage.Reservation) (*storage.Reservation, error) {
	err := r.DB.DB.Create(reservation).Error
	if err != nil {
		return nil, err
	}
	// Reload with preloads
	return r.FindByID(reservation.ID)
}

// Update сохраняет изменения бронирования и возвращает обновлённую версию с предзагрузкой.
func (r *reservationRepository) Update(reservation *storage.Reservation) (*storage.Reservation, error) {
	err := r.DB.DB.Save(reservation).Error
	if err != nil {
		return nil, err
	}
	// Reload with preloads
	return r.FindByID(reservation.ID)
}

// Delete полностью удаляет бронирование из базы данных.
func (r *reservationRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.Reservation{}, id).Error
}

// UpdateStatus обновляет только статус бронирования.
// Используется для отмены бронирования (Cancel, CancelByPhone).
func (r *reservationRepository) UpdateStatus(id uint, status storage.ReservationStatus) error {
	return r.DB.DB.Model(&storage.Reservation{}).Where("id = ?", id).Update("status", status).Error
}
