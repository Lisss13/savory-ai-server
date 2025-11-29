// Package service содержит бизнес-логику для модуля поддержки.
package service

import (
	"errors"
	"savory-ai-server/app/module/support/payload"
	"savory-ai-server/app/module/support/repository"
	"savory-ai-server/app/storage"
)

type supportService struct {
	repo repository.SupportRepository
}

// SupportService определяет интерфейс бизнес-логики для работы с заявками в поддержку.
type SupportService interface {
	// =====================================================
	// User Operations
	// =====================================================

	// Create создаёт новую заявку в поддержку
	// Метод: POST /support
	Create(userID uint, req *payload.CreateSupportTicketReq) (*payload.SupportTicketResp, error)

	// GetMyTickets возвращает заявки текущего пользователя
	// Метод: GET /support/my
	GetMyTickets(userID uint, page, pageSize int) (*payload.SupportTicketsListResp, error)

	// GetByID возвращает заявку по ID
	// Метод: GET /support/:id
	GetByID(id uint) (*payload.SupportTicketResp, error)

	// =====================================================
	// Admin Operations
	// =====================================================

	// GetAll возвращает все заявки (для админов)
	// Метод: GET /admin/support
	GetAll(page, pageSize int) (*payload.SupportTicketsListResp, error)

	// GetByStatus возвращает заявки по статусу (для админов)
	// Метод: GET /admin/support?status=in_progress
	GetByStatus(status string, page, pageSize int) (*payload.SupportTicketsListResp, error)

	// UpdateStatus обновляет статус заявки (для админов)
	// Метод: PATCH /admin/support/:id/status
	UpdateStatus(id uint, req *payload.UpdateSupportTicketStatusReq) (*payload.SupportTicketResp, error)
}

// NewSupportService создаёт новый экземпляр сервиса поддержки
func NewSupportService(repo repository.SupportRepository) SupportService {
	return &supportService{repo: repo}
}

// Create создаёт новую заявку в поддержку
func (s *supportService) Create(userID uint, req *payload.CreateSupportTicketReq) (*payload.SupportTicketResp, error) {
	ticket := &storage.SupportTicket{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Email:       req.Email,
		Phone:       req.Phone,
		Status:      storage.SupportTicketStatusInProgress, // По умолчанию "взят в работу"
	}

	created, err := s.repo.Create(ticket)
	if err != nil {
		return nil, err
	}

	return s.ticketToResp(created), nil
}

// GetMyTickets возвращает заявки текущего пользователя
func (s *supportService) GetMyTickets(userID uint, page, pageSize int) (*payload.SupportTicketsListResp, error) {
	tickets, total, err := s.repo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return s.ticketsToListResp(tickets, total, page, pageSize), nil
}

// GetByID возвращает заявку по ID
func (s *supportService) GetByID(id uint) (*payload.SupportTicketResp, error) {
	ticket, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("заявка не найдена")
	}

	return s.ticketToResp(ticket), nil
}

// GetAll возвращает все заявки (для админов)
func (s *supportService) GetAll(page, pageSize int) (*payload.SupportTicketsListResp, error) {
	tickets, total, err := s.repo.FindAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	return s.ticketsToListResp(tickets, total, page, pageSize), nil
}

// GetByStatus возвращает заявки по статусу (для админов)
func (s *supportService) GetByStatus(status string, page, pageSize int) (*payload.SupportTicketsListResp, error) {
	ticketStatus := storage.SupportTicketStatus(status)
	if ticketStatus != storage.SupportTicketStatusInProgress && ticketStatus != storage.SupportTicketStatusCompleted {
		return nil, errors.New("неверный статус: допустимые значения in_progress, completed")
	}

	tickets, total, err := s.repo.FindByStatus(ticketStatus, page, pageSize)
	if err != nil {
		return nil, err
	}

	return s.ticketsToListResp(tickets, total, page, pageSize), nil
}

// UpdateStatus обновляет статус заявки (для админов)
func (s *supportService) UpdateStatus(id uint, req *payload.UpdateSupportTicketStatusReq) (*payload.SupportTicketResp, error) {
	// Проверяем существование заявки
	_, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("заявка не найдена")
	}

	// Обновляем статус
	newStatus := storage.SupportTicketStatus(req.Status)
	if err := s.repo.UpdateStatus(id, newStatus); err != nil {
		return nil, err
	}

	// Возвращаем обновлённую заявку
	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.ticketToResp(updated), nil
}

// =====================================================
// Helpers
// =====================================================

// ticketToResp преобразует модель заявки в ответ
func (s *supportService) ticketToResp(ticket *storage.SupportTicket) *payload.SupportTicketResp {
	return &payload.SupportTicketResp{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		UserName:    ticket.User.Name,
		UserEmail:   ticket.User.Email,
		Title:       ticket.Title,
		Description: ticket.Description,
		Email:       ticket.Email,
		Phone:       ticket.Phone,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}
}

// ticketsToListResp преобразует список заявок в ответ с пагинацией
func (s *supportService) ticketsToListResp(tickets []storage.SupportTicket, total int64, page, pageSize int) *payload.SupportTicketsListResp {
	var ticketResps []payload.SupportTicketResp
	for _, ticket := range tickets {
		ticketResps = append(ticketResps, *s.ticketToResp(&ticket))
	}

	return &payload.SupportTicketsListResp{
		Tickets:    ticketResps,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}
}
