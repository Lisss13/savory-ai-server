// Package repository содержит слой доступа к данным для модуля чата.
// Использует GORM для работы с PostgreSQL.
package repository

import (
	"gorm.io/gorm"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
	"time"
)

// chatRepository реализует интерфейс ChatRepository.
type chatRepository struct {
	DB *database.Database
}

// ChatRepository определяет интерфейс для работы с чат-сессиями и сообщениями в БД.
// Разделён на две группы: Table Chat и Restaurant Chat.
type ChatRepository interface {
	// =====================================================
	// Table Chat - чат для посетителей за столиком
	// =====================================================
	FindTableSessionByRestaurantID(restaurantID uint) ([]*storage.TableChatSessions, error) // Найти сессии ресторана
	FindTableSessionByID(sessionID uint) (*storage.TableChatSessions, error)                // Найти сессию по ID
	CreateTableSession(session *storage.TableChatSessions) (*storage.TableChatSessions, error) // Создать сессию
	CloseTableSession(sessionID uint) error                                                  // Закрыть сессию
	CreateTableMessage(message *storage.TableChatMessage) (*storage.TableChatMessage, error) // Создать сообщение
	UpdateTableSessionActivity(sessionID uint) error                                         // Обновить время активности
	GetMessagesBySessionID(sessionID uint) ([]*storage.TableChatMessage, error)             // Получить сообщения сессии
	GetSessionsByTableID(tableID uint) ([]*storage.TableChatSessions, error)                // Получить сессии столика

	// =====================================================
	// Restaurant Chat - общий чат с рестораном
	// =====================================================
	FindRestaurantSessionByID(sessionID uint) (*storage.RestaurantChatSessions, error)                   // Найти сессию по ID
	FindRestaurantSessionByRestaurantID(restaurantID uint) ([]*storage.RestaurantChatSessions, error)    // Найти сессии ресторана
	CreateRestaurantSession(session *storage.RestaurantChatSessions) (*storage.RestaurantChatSessions, error) // Создать сессию
	CloseRestaurantSession(sessionID uint) error                                                         // Закрыть сессию
	CreateRestaurantMessage(message *storage.RestaurantChatMessage) (*storage.RestaurantChatMessage, error) // Создать сообщение
	UpdateRestaurantSessionActivity(sessionID uint) error                                                // Обновить время активности
	GetRestaurantMessagesBySessionID(sessionID uint) ([]*storage.RestaurantChatMessage, error)           // Получить сообщения
}

// NewChatRepository создаёт новый экземпляр репозитория чата.
func NewChatRepository(db *database.Database) ChatRepository {
	return &chatRepository{
		DB: db,
	}
}

// =====================================================
// Table Chat Operations
// =====================================================

// FindTableSessionByRestaurantID находит все чат-сессии столиков для ресторана.
// Предзагружает связи: Table, Restaurant, Messages.
func (r *chatRepository) FindTableSessionByRestaurantID(restaurantID uint) ([]*storage.TableChatSessions, error) {
	var sessions []*storage.TableChatSessions

	res := r.DB.DB.
		Preload("Table").
		Preload("Restaurant").
		Preload("Messages").
		Where("restaurant_id = ?", restaurantID).
		Find(&sessions)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return sessions, nil
}

// CreateTableSession создаёт новую чат-сессию для столика.
// Возвращает сессию с предзагруженной связью Table.
func (r *chatRepository) CreateTableSession(session *storage.TableChatSessions) (*storage.TableChatSessions, error) {
	if err := r.DB.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	// Reload the session with associations
	var createdSession storage.TableChatSessions
	if err := r.DB.DB.Preload("Table").First(&createdSession, session.ID).Error; err != nil {
		return nil, err
	}

	return &createdSession, nil
}

// CloseTableSession закрывает сессию чата (устанавливает active = false).
func (r *chatRepository) CloseTableSession(sessionID uint) error {
	return r.DB.DB.Model(&storage.TableChatSessions{}).
		Where("id = ?", sessionID).
		Update("active", false).Error
}

// FindTableSessionByID находит сессию чата столика по ID.
// Предзагружает связь Table.
func (r *chatRepository) FindTableSessionByID(sessionID uint) (*storage.TableChatSessions, error) {
	var session storage.TableChatSessions

	err := r.DB.DB.
		Preload("Table").
		Where("id = ?", sessionID).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// CreateTableMessage создаёт сообщение в чате столика.
// Автоматически устанавливает sent_at если не указано.
// Возвращает сообщение с предзагруженной связью Table.
func (r *chatRepository) CreateTableMessage(message *storage.TableChatMessage) (*storage.TableChatMessage, error) {
	// Set sent_at to current time if not provided
	if message.SentAt.IsZero() {
		message.SentAt = time.Now()
	}

	if err := r.DB.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	// Reload the message with associations
	var createdMessage storage.TableChatMessage
	if err := r.DB.DB.
		Preload("Table").
		First(&createdMessage, message.ID).Error; err != nil {
		return nil, err
	}

	return &createdMessage, nil
}

// UpdateTableSessionActivity обновляет время последней активности сессии.
// Также устанавливает active = true.
func (r *chatRepository) UpdateTableSessionActivity(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active": time.Now(),
			"active":      true,
		}).Error
}

// GetMessagesBySessionID возвращает все сообщения из сессии чата столика.
// Сортировка по времени отправки (старые сначала).
func (r *chatRepository) GetMessagesBySessionID(sessionID uint) ([]*storage.TableChatMessage, error) {
	var messages []*storage.TableChatMessage

	query := r.DB.DB.
		Preload("Table").
		Where("chat_session_id = ?", sessionID).
		Order("sent_at ASC")

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// GetSessionsByTableID возвращает все чат-сессии для столика.
// Предзагружает связи: Table, Messages.
// Сортировка по времени последней активности.
func (r *chatRepository) GetSessionsByTableID(tableID uint) ([]*storage.TableChatSessions, error) {
	var sessions []*storage.TableChatSessions

	query := r.DB.DB.
		Preload("Table").
		Preload("Messages").
		Where("table_id = ?", tableID).
		Order("last_active ASC")

	if err := query.Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// =====================================================
// Restaurant Chat Operations
// =====================================================

// FindRestaurantSessionByID находит сессию чата ресторана по ID.
// Предзагружает связь Restaurant.
func (r *chatRepository) FindRestaurantSessionByID(sessionID uint) (*storage.RestaurantChatSessions, error) {
	var session storage.RestaurantChatSessions

	err := r.DB.DB.
		Preload("Restaurant").
		Where("id = ?", sessionID).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// FindRestaurantSessionByRestaurantID находит все чат-сессии ресторана.
// Предзагружает связи: Restaurant, Messages.
// Возвращает gorm.ErrRecordNotFound если сессий нет.
func (r *chatRepository) FindRestaurantSessionByRestaurantID(restaurantID uint) ([]*storage.RestaurantChatSessions, error) {
	var sessions []*storage.RestaurantChatSessions

	res := r.DB.DB.
		Preload("Restaurant").
		Preload("Messages").
		Where("restaurant_id = ?", restaurantID).
		Find(&sessions)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return sessions, nil
}

// CreateRestaurantSession создаёт новую сессию чата для ресторана.
// Возвращает сессию с предзагруженной связью Restaurant.
func (r *chatRepository) CreateRestaurantSession(session *storage.RestaurantChatSessions) (*storage.RestaurantChatSessions, error) {
	if err := r.DB.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	// Reload the session with associations
	var createdSession storage.RestaurantChatSessions
	if err := r.DB.DB.Preload("Restaurant").First(&createdSession, session.ID).Error; err != nil {
		return nil, err
	}

	return &createdSession, nil
}

// CloseRestaurantSession закрывает сессию чата ресторана (устанавливает active = false).
func (r *chatRepository) CloseRestaurantSession(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Update("active", false).Error
}

// CreateRestaurantMessage создаёт сообщение в чате ресторана.
// Автоматически устанавливает sent_at если не указано.
// Возвращает сообщение с предзагруженной связью Restaurant.
func (r *chatRepository) CreateRestaurantMessage(message *storage.RestaurantChatMessage) (*storage.RestaurantChatMessage, error) {
	// Set sent_at to current time if not provided
	if message.SentAt.IsZero() {
		message.SentAt = time.Now()
	}

	if err := r.DB.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	// Reload the message with associations
	var createdMessage storage.RestaurantChatMessage
	if err := r.DB.DB.
		Preload("Restaurant").
		First(&createdMessage, message.ID).Error; err != nil {
		return nil, err
	}

	return &createdMessage, nil
}

// UpdateRestaurantSessionActivity обновляет время последней активности сессии ресторана.
// Также устанавливает active = true.
func (r *chatRepository) UpdateRestaurantSessionActivity(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active": time.Now(),
			"active":      true,
		}).Error
}

// GetRestaurantMessagesBySessionID возвращает все сообщения из сессии чата ресторана.
// Предзагружает связь Restaurant.
// Сортировка по времени отправки (старые сначала).
func (r *chatRepository) GetRestaurantMessagesBySessionID(sessionID uint) ([]*storage.RestaurantChatMessage, error) {
	var messages []*storage.RestaurantChatMessage

	query := r.DB.DB.
		Preload("Restaurant").
		Where("chat_session_id = ?", sessionID).
		Order("sent_at ASC")

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
