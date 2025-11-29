// Package repository содержит слой доступа к данным для модуля поддержки.
// Использует GORM для работы с PostgreSQL.
package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type supportRepository struct {
	DB *database.Database
}

// SupportRepository определяет интерфейс для работы с заявками в поддержку.
type SupportRepository interface {
	// =====================================================
	// CRUD Operations
	// =====================================================

	// Create создаёт новую заявку в поддержку
	Create(ticket *storage.SupportTicket) (*storage.SupportTicket, error)

	// FindByID находит заявку по ID
	FindByID(id uint) (*storage.SupportTicket, error)

	// FindByUserID находит все заявки пользователя
	FindByUserID(userID uint, page, pageSize int) ([]storage.SupportTicket, int64, error)

	// FindAll находит все заявки (для админов)
	FindAll(page, pageSize int) ([]storage.SupportTicket, int64, error)

	// FindByStatus находит заявки по статусу (для админов)
	FindByStatus(status storage.SupportTicketStatus, page, pageSize int) ([]storage.SupportTicket, int64, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(id uint, status storage.SupportTicketStatus) error
}

// NewSupportRepository создаёт новый экземпляр репозитория поддержки
func NewSupportRepository(db *database.Database) SupportRepository {
	return &supportRepository{DB: db}
}

// Create создаёт новую заявку в поддержку
func (r *supportRepository) Create(ticket *storage.SupportTicket) (*storage.SupportTicket, error) {
	err := r.DB.DB.Create(ticket).Error
	if err != nil {
		return nil, err
	}
	// Загружаем пользователя для ответа
	r.DB.DB.Preload("User").First(ticket, ticket.ID)
	return ticket, nil
}

// FindByID находит заявку по ID с данными пользователя
func (r *supportRepository) FindByID(id uint) (*storage.SupportTicket, error) {
	var ticket storage.SupportTicket
	err := r.DB.DB.Preload("User").First(&ticket, id).Error
	return &ticket, err
}

// FindByUserID находит все заявки пользователя с пагинацией
func (r *supportRepository) FindByUserID(userID uint, page, pageSize int) ([]storage.SupportTicket, int64, error) {
	var tickets []storage.SupportTicket
	var total int64

	offset := (page - 1) * pageSize

	// Считаем общее количество
	r.DB.DB.Model(&storage.SupportTicket{}).Where("user_id = ?", userID).Count(&total)

	// Получаем заявки с пагинацией
	err := r.DB.DB.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tickets).Error

	return tickets, total, err
}

// FindAll находит все заявки с пагинацией (для админов)
func (r *supportRepository) FindAll(page, pageSize int) ([]storage.SupportTicket, int64, error) {
	var tickets []storage.SupportTicket
	var total int64

	offset := (page - 1) * pageSize

	// Считаем общее количество
	r.DB.DB.Model(&storage.SupportTicket{}).Count(&total)

	// Получаем заявки с пагинацией
	err := r.DB.DB.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tickets).Error

	return tickets, total, err
}

// FindByStatus находит заявки по статусу с пагинацией (для админов)
func (r *supportRepository) FindByStatus(status storage.SupportTicketStatus, page, pageSize int) ([]storage.SupportTicket, int64, error) {
	var tickets []storage.SupportTicket
	var total int64

	offset := (page - 1) * pageSize

	// Считаем общее количество
	r.DB.DB.Model(&storage.SupportTicket{}).Where("status = ?", status).Count(&total)

	// Получаем заявки с пагинацией
	err := r.DB.DB.Preload("User").
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tickets).Error

	return tickets, total, err
}

// UpdateStatus обновляет статус заявки
func (r *supportRepository) UpdateStatus(id uint, status storage.SupportTicketStatus) error {
	return r.DB.DB.Model(&storage.SupportTicket{}).Where("id = ?", id).Update("status", status).Error
}
