package service

import (
	"encoding/json"
	"savory-ai-server/app/module/admin/payload"
	"savory-ai-server/app/module/admin/repository"
	"savory-ai-server/app/storage"
)

type adminService struct {
	adminRepo repository.AdminRepository
}

type AdminService interface {
	// Статистика
	GetStats() (*payload.StatsResp, error)

	// Управление пользователями
	GetAllUsers(page, pageSize int) (*payload.AdminUsersResp, error)
	GetUserByID(id uint) (*payload.AdminUserResp, error)
	UpdateUserStatus(adminID uint, userID uint, isActive bool, ipAddress string) error
	UpdateUserRole(adminID uint, userID uint, role string, ipAddress string) error
	DeleteUser(adminID uint, userID uint, ipAddress string) error

	// Управление организациями
	GetAllOrganizations(page, pageSize int) (*payload.AdminOrganizationsResp, error)
	GetOrganizationByID(id uint) (*payload.AdminOrganizationResp, error)
	DeleteOrganization(adminID uint, orgID uint, ipAddress string) error

	// Модерация контента
	GetAllDishes(page, pageSize int) (*payload.AdminDishesResp, error)
	DeleteDish(adminID uint, dishID uint, ipAddress string) error

	// Логи
	GetAllLogs(page, pageSize int) (*payload.AdminLogsResp, error)
	GetLogsByAdminID(adminID uint, page, pageSize int) (*payload.AdminLogsResp, error)
}

func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{adminRepo: adminRepo}
}

// ==================== Статистика ====================

// GetStats - Получение статистики системы
//
// Возвращает общую статистику:
// - Количество пользователей (всего и активных)
// - Количество организаций, ресторанов, блюд, столиков
// - Количество активных подписок
func (s *adminService) GetStats() (*payload.StatsResp, error) {
	totalUsers, _ := s.adminRepo.CountUsers()
	activeUsers, _ := s.adminRepo.CountActiveUsers()
	totalOrgs, _ := s.adminRepo.CountOrganizations()
	totalRestaurants, _ := s.adminRepo.CountRestaurants()
	totalDishes, _ := s.adminRepo.CountDishes()
	totalTables, _ := s.adminRepo.CountTables()
	totalQuestions, _ := s.adminRepo.CountQuestions()
	activeSubscriptions, _ := s.adminRepo.CountActiveSubscriptions()

	// Получаем недавнюю активность
	logs, _, _ := s.adminRepo.FindAllLogs(1, 10)
	var recentActivity []payload.ActivityResp
	for _, log := range logs {
		recentActivity = append(recentActivity, payload.ActivityResp{
			ID:         log.ID,
			Action:     string(log.Action),
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			AdminName:  log.Admin.Name,
			CreatedAt:  log.CreatedAt,
		})
	}

	return &payload.StatsResp{
		TotalUsers:          totalUsers,
		ActiveUsers:         activeUsers,
		TotalOrganizations:  totalOrgs,
		TotalRestaurants:    totalRestaurants,
		TotalDishes:         totalDishes,
		TotalTables:         totalTables,
		TotalQuestions:      totalQuestions,
		ActiveSubscriptions: activeSubscriptions,
		RecentActivity:      recentActivity,
	}, nil
}

// ==================== Пользователи ====================

// GetAllUsers - Получение списка всех пользователей с пагинацией
func (s *adminService) GetAllUsers(page, pageSize int) (*payload.AdminUsersResp, error) {
	users, total, err := s.adminRepo.FindAllUsers(page, pageSize)
	if err != nil {
		return nil, err
	}

	var userResps []payload.AdminUserResp
	for _, user := range users {
		userResps = append(userResps, payload.AdminUserResp{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Company:   user.Company,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		})
	}

	return &payload.AdminUsersResp{
		Users:      userResps,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetUserByID - Получение пользователя по ID
func (s *adminService) GetUserByID(id uint) (*payload.AdminUserResp, error) {
	user, err := s.adminRepo.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return &payload.AdminUserResp{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Company:   user.Company,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateUserStatus - Блокировка/разблокировка пользователя
func (s *adminService) UpdateUserStatus(adminID uint, userID uint, isActive bool, ipAddress string) error {
	if err := s.adminRepo.UpdateUserStatus(userID, isActive); err != nil {
		return err
	}

	// Логируем действие
	action := storage.ActionBlock
	if isActive {
		action = storage.ActionUnblock
	}
	s.logAction(adminID, action, "user", userID, map[string]interface{}{"isActive": isActive}, ipAddress)

	return nil
}

// UpdateUserRole - Изменение роли пользователя
func (s *adminService) UpdateUserRole(adminID uint, userID uint, role string, ipAddress string) error {
	if err := s.adminRepo.UpdateUserRole(userID, storage.UserRole(role)); err != nil {
		return err
	}

	s.logAction(adminID, storage.ActionUpdate, "user", userID, map[string]interface{}{"role": role}, ipAddress)
	return nil
}

// DeleteUser - Удаление пользователя
func (s *adminService) DeleteUser(adminID uint, userID uint, ipAddress string) error {
	if err := s.adminRepo.DeleteUser(userID); err != nil {
		return err
	}

	s.logAction(adminID, storage.ActionDelete, "user", userID, nil, ipAddress)
	return nil
}

// ==================== Организации ====================

// GetAllOrganizations - Получение списка всех организаций
func (s *adminService) GetAllOrganizations(page, pageSize int) (*payload.AdminOrganizationsResp, error) {
	orgs, total, err := s.adminRepo.FindAllOrganizations(page, pageSize)
	if err != nil {
		return nil, err
	}

	var orgResps []payload.AdminOrganizationResp
	for _, org := range orgs {
		orgResps = append(orgResps, payload.AdminOrganizationResp{
			ID:         org.ID,
			Name:       org.Name,
			Phone:      org.Phone,
			AdminID:    org.AdminID,
			AdminName:  org.Admin.Name,
			AdminEmail: org.Admin.Email,
			CreatedAt:  org.CreatedAt,
		})
	}

	return &payload.AdminOrganizationsResp{
		Organizations: orgResps,
		TotalCount:    total,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

// GetOrganizationByID - Получение организации по ID
func (s *adminService) GetOrganizationByID(id uint) (*payload.AdminOrganizationResp, error) {
	org, err := s.adminRepo.FindOrganizationByID(id)
	if err != nil {
		return nil, err
	}

	return &payload.AdminOrganizationResp{
		ID:         org.ID,
		Name:       org.Name,
		Phone:      org.Phone,
		AdminID:    org.AdminID,
		AdminName:  org.Admin.Name,
		AdminEmail: org.Admin.Email,
		CreatedAt:  org.CreatedAt,
	}, nil
}

// DeleteOrganization - Удаление организации
func (s *adminService) DeleteOrganization(adminID uint, orgID uint, ipAddress string) error {
	if err := s.adminRepo.DeleteOrganization(orgID); err != nil {
		return err
	}

	s.logAction(adminID, storage.ActionDelete, "organization", orgID, nil, ipAddress)
	return nil
}

// ==================== Модерация контента ====================

// GetAllDishes - Получение всех блюд для модерации
func (s *adminService) GetAllDishes(page, pageSize int) (*payload.AdminDishesResp, error) {
	dishes, total, err := s.adminRepo.FindAllDishes(page, pageSize)
	if err != nil {
		return nil, err
	}

	var dishResps []payload.AdminDishResp
	for _, dish := range dishes {
		dishResps = append(dishResps, payload.AdminDishResp{
			ID:               dish.ID,
			Name:             dish.Name,
			Description:      dish.Description,
			Price:            dish.Price,
			Image:            dish.Image,
			OrganizationID:   dish.OrganizationID,
			OrganizationName: dish.Organization.Name,
			CreatedAt:        dish.CreatedAt,
		})
	}

	return &payload.AdminDishesResp{
		Dishes:     dishResps,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// DeleteDish - Удаление блюда (модерация)
func (s *adminService) DeleteDish(adminID uint, dishID uint, ipAddress string) error {
	if err := s.adminRepo.DeleteDish(dishID); err != nil {
		return err
	}

	s.logAction(adminID, storage.ActionDelete, "dish", dishID, nil, ipAddress)
	return nil
}

// ==================== Логи ====================

// GetAllLogs - Получение всех логов действий администраторов
func (s *adminService) GetAllLogs(page, pageSize int) (*payload.AdminLogsResp, error) {
	logs, total, err := s.adminRepo.FindAllLogs(page, pageSize)
	if err != nil {
		return nil, err
	}

	var logResps []payload.AdminLogResp
	for _, log := range logs {
		logResps = append(logResps, payload.AdminLogResp{
			ID:         log.ID,
			AdminID:    log.AdminID,
			AdminName:  log.Admin.Name,
			AdminEmail: log.Admin.Email,
			Action:     string(log.Action),
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			Details:    log.Details,
			IPAddress:  log.IPAddress,
			CreatedAt:  log.CreatedAt,
		})
	}

	return &payload.AdminLogsResp{
		Logs:       logResps,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetLogsByAdminID - Получение логов конкретного администратора
func (s *adminService) GetLogsByAdminID(adminID uint, page, pageSize int) (*payload.AdminLogsResp, error) {
	logs, total, err := s.adminRepo.FindLogsByAdminID(adminID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var logResps []payload.AdminLogResp
	for _, log := range logs {
		logResps = append(logResps, payload.AdminLogResp{
			ID:         log.ID,
			AdminID:    log.AdminID,
			AdminName:  log.Admin.Name,
			AdminEmail: log.Admin.Email,
			Action:     string(log.Action),
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			Details:    log.Details,
			IPAddress:  log.IPAddress,
			CreatedAt:  log.CreatedAt,
		})
	}

	return &payload.AdminLogsResp{
		Logs:       logResps,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// ==================== Helpers ====================

// logAction - вспомогательная функция для логирования действий администратора
func (s *adminService) logAction(adminID uint, action storage.AdminAction, entityType string, entityID uint, details map[string]interface{}, ipAddress string) {
	detailsJSON := ""
	if details != nil {
		if bytes, err := json.Marshal(details); err == nil {
			detailsJSON = string(bytes)
		}
	}

	log := &storage.AdminLog{
		AdminID:    adminID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    detailsJSON,
		IPAddress:  ipAddress,
	}

	s.adminRepo.CreateLog(log)
}
