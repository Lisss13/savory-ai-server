package service

import (
	"errors"
	"savory-ai-server/app/module/subscription/payload"
	"savory-ai-server/app/module/subscription/repository"
	"savory-ai-server/app/storage"
	"time"
)

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
}

type SubscriptionService interface {
	GetAll() (*payload.SubscriptionsResp, error)
	GetByID(id uint) (*payload.SubscriptionResp, error)
	GetByOrganizationID(organizationID uint) (*payload.SubscriptionsResp, error)
	GetActiveByOrganizationID(organizationID uint) (*payload.SubscriptionResp, error)
	Create(req *payload.CreateSubscriptionReq) (*payload.SubscriptionResp, error)
	Update(id uint, req *payload.UpdateSubscriptionReq) (*payload.SubscriptionResp, error)
	Extend(id uint, req *payload.ExtendSubscriptionReq) (*payload.SubscriptionResp, error)
	Deactivate(id uint) (*payload.SubscriptionResp, error)
	Delete(id uint) error
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

// GetAll - Получение всех подписок в системе
//
// Бизнес-логика:
// - Загружает все подписки из базы данных
// - Преобразует каждую подписку в формат ответа API
// - Используется администратором для просмотра всех подписок
func (s *subscriptionService) GetAll() (*payload.SubscriptionsResp, error) {
	subscriptions, err := s.subscriptionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var subscriptionResps []payload.SubscriptionResp
	for _, subscription := range subscriptions {
		subscriptionResps = append(subscriptionResps, mapSubscriptionToResponse(subscription))
	}

	return &payload.SubscriptionsResp{
		Subscriptions: subscriptionResps,
	}, nil
}

// GetByID - Получение подписки по её идентификатору
//
// Бизнес-логика:
// - Ищет подписку по ID в базе данных
// - Возвращает ошибку, если подписка не найдена
// - Используется для просмотра детальной информации о конкретной подписке
func (s *subscriptionService) GetByID(id uint) (*payload.SubscriptionResp, error) {
	subscription, err := s.subscriptionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(subscription)
	return &resp, nil
}

// GetByOrganizationID - Получение всех подписок организации
//
// Бизнес-логика:
// - Загружает все подписки (активные и неактивные) для указанной организации
// - Сортирует по дате создания (новые первыми)
// - Позволяет просмотреть историю подписок организации
func (s *subscriptionService) GetByOrganizationID(organizationID uint) (*payload.SubscriptionsResp, error) {
	subscriptions, err := s.subscriptionRepo.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var subscriptionResps []payload.SubscriptionResp
	for _, subscription := range subscriptions {
		subscriptionResps = append(subscriptionResps, mapSubscriptionToResponse(subscription))
	}

	return &payload.SubscriptionsResp{
		Subscriptions: subscriptionResps,
	}, nil
}

// GetActiveByOrganizationID - Получение активной подписки организации
//
// Бизнес-логика:
// - Ищет текущую активную подписку для организации
// - У организации может быть только одна активная подписка
// - Возвращает ошибку если активной подписки нет
// - Используется для проверки доступа организации к системе
func (s *subscriptionService) GetActiveByOrganizationID(organizationID uint) (*payload.SubscriptionResp, error) {
	subscription, err := s.subscriptionRepo.FindActiveByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(subscription)
	return &resp, nil
}

// Create - Создание новой подписки для организации
//
// Бизнес-логика:
// 1. Проверяет наличие активной подписки у организации
// 2. Если есть активная подписка - деактивирует её (у организации только одна активная)
// 3. Вычисляет дату окончания: StartDate + Period месяцев
// 4. Создает новую подписку со статусом IsActive = true
//
// Пример: если Period = 6 и StartDate = 1 января,
// то EndDate будет 1 июля (через 6 месяцев)
func (s *subscriptionService) Create(req *payload.CreateSubscriptionReq) (*payload.SubscriptionResp, error) {
	// Деактивируем существующую активную подписку организации
	existingActive, err := s.subscriptionRepo.FindActiveByOrganizationID(req.OrganizationID)
	if err == nil && existingActive != nil {
		existingActive.IsActive = false
		if _, err := s.subscriptionRepo.Update(existingActive); err != nil {
			return nil, err
		}
	}

	// Вычисляем дату окончания на основе периода (в месяцах)
	endDate := req.StartDate.AddDate(0, req.Period, 0)

	subscription := &storage.Subscription{
		OrganizationID: req.OrganizationID,
		Period:         req.Period,
		StartDate:      req.StartDate,
		EndDate:        endDate,
		IsActive:       true,
	}

	createdSubscription, err := s.subscriptionRepo.Create(subscription)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(createdSubscription)
	return &resp, nil
}

// Update - Обновление существующей подписки
//
// Бизнес-логика:
// 1. Проверяет существование подписки по ID
// 2. Позволяет изменить: период, дату начала, статус активности
// 3. Автоматически пересчитывает дату окончания при изменении периода или даты начала
//
// Используется для:
// - Корректировки ошибочно созданных подписок
// - Изменения условий подписки
// - Ручной активации/деактивации
func (s *subscriptionService) Update(id uint, req *payload.UpdateSubscriptionReq) (*payload.SubscriptionResp, error) {
	subscription, err := s.subscriptionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Пересчитываем дату окончания на основе нового периода
	endDate := req.StartDate.AddDate(0, req.Period, 0)

	subscription.Period = req.Period
	subscription.StartDate = req.StartDate
	subscription.EndDate = endDate
	subscription.IsActive = req.IsActive

	updatedSubscription, err := s.subscriptionRepo.Update(subscription)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(updatedSubscription)
	return &resp, nil
}

// Extend - Продление подписки на дополнительные месяцы
//
// Бизнес-логика:
// 1. Проверяет существование подписки
// 2. Проверяет что подписка активна (нельзя продлить неактивную)
// 3. Добавляет указанное количество месяцев к текущему периоду
// 4. Сдвигает дату окончания на соответствующее количество месяцев
//
// Пример: подписка на 6 месяцев до 1 июля, продление на 3 месяца
// Результат: период = 9, дата окончания = 1 октября
//
// Ошибка: если подписка неактивна - нужно создать новую
func (s *subscriptionService) Extend(id uint, req *payload.ExtendSubscriptionReq) (*payload.SubscriptionResp, error) {
	subscription, err := s.subscriptionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Нельзя продлить неактивную подписку
	if !subscription.IsActive {
		return nil, errors.New("cannot extend inactive subscription")
	}

	// Добавляем месяцы к периоду и сдвигаем дату окончания
	subscription.Period += req.Period
	subscription.EndDate = subscription.EndDate.AddDate(0, req.Period, 0)

	updatedSubscription, err := s.subscriptionRepo.Update(subscription)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(updatedSubscription)
	return &resp, nil
}

// Deactivate - Деактивация подписки
//
// Бизнес-логика:
// - Устанавливает IsActive = false для указанной подписки
// - Подписка остается в базе данных для истории
// - Организация теряет доступ к функционалу системы
//
// Используется, когда:
// - Организация отказалась от подписки
// - Нужно временно заблокировать доступ
// - Организация нарушила условия использования
func (s *subscriptionService) Deactivate(id uint) (*payload.SubscriptionResp, error) {
	subscription, err := s.subscriptionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	subscription.IsActive = false

	updatedSubscription, err := s.subscriptionRepo.Update(subscription)
	if err != nil {
		return nil, err
	}

	resp := mapSubscriptionToResponse(updatedSubscription)
	return &resp, nil
}

// Delete - Удаление подписки из системы
//
// Бизнес-логика:
// - Полностью удаляет запись о подписке (soft delete через GORM)
// - Используется для удаления ошибочно созданных подписок
// - Рекомендуется использовать Deactivate вместо Delete для сохранения истории
func (s *subscriptionService) Delete(id uint) error {
	return s.subscriptionRepo.Delete(id)
}

// mapSubscriptionToResponse - Преобразование модели подписки в формат ответа API
//
// Бизнес-логика:
// 1. Вычисляет количество оставшихся дней до окончания подписки
// 2. Если подписка истекла - DaysLeft = 0
// 3. Преобразует связанную организацию в формат ответа
// 4. Формирует полный ответ со всеми полями
//
// DaysLeft используется клиентом для:
// - Отображения предупреждений о скором окончании
// - Напоминаний о необходимости продления
func mapSubscriptionToResponse(subscription *storage.Subscription) payload.SubscriptionResp {
	// Вычисляем количество оставшихся дней
	daysLeft := int(time.Until(subscription.EndDate).Hours() / 24)
	if daysLeft < 0 {
		daysLeft = 0
	}

	// Преобразуем организацию
	organizationResp := payload.OrganizationResp{
		ID:    subscription.Organization.ID,
		Name:  subscription.Organization.Name,
		Phone: subscription.Organization.Phone,
	}

	return payload.SubscriptionResp{
		ID:           subscription.ID,
		CreatedAt:    subscription.CreatedAt,
		Organization: organizationResp,
		Period:       subscription.Period,
		StartDate:    subscription.StartDate,
		EndDate:      subscription.EndDate,
		IsActive:     subscription.IsActive,
		DaysLeft:     daysLeft,
	}
}
