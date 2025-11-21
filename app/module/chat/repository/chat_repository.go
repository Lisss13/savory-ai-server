package repository

import (
	"gorm.io/gorm"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
	"time"
)

type chatRepository struct {
	DB *database.Database
}

type ChatRepository interface {
	// Table chat operations
	FindTableSessionByRestaurantID(restaurantID uint) ([]*storage.TableChatSessions, error)
	CreateTableSession(session *storage.TableChatSessions) (*storage.TableChatSessions, error)
	CloseTableSession(sessionID uint) error
	FindTableSessionByID(sessionID uint) (*storage.TableChatSessions, error)
	CreateTableMessage(message *storage.TableChatMessage) (*storage.TableChatMessage, error)
	UpdateTableSessionActivity(sessionID uint) error
	GetMessagesBySessionID(sessionID uint) ([]*storage.TableChatMessage, error)
	GetSessionsByTableID(tableID uint) ([]*storage.TableChatSessions, error)

	// Restaurant chat operations
	FindRestaurantSessionByID(sessionID uint) (*storage.RestaurantChatSessions, error)
	FindRestaurantSessionByRestaurantID(restaurantID uint) ([]*storage.RestaurantChatSessions, error)
	CreateRestaurantSession(session *storage.RestaurantChatSessions) (*storage.RestaurantChatSessions, error)
	CloseRestaurantSession(sessionID uint) error
	CreateRestaurantMessage(message *storage.RestaurantChatMessage) (*storage.RestaurantChatMessage, error)
	UpdateRestaurantSessionActivity(sessionID uint) error
	GetRestaurantMessagesBySessionID(sessionID uint) ([]*storage.RestaurantChatMessage, error)
}

func NewChatRepository(db *database.Database) ChatRepository {
	return &chatRepository{
		DB: db,
	}
}

// ChatMessage operations

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

// CreateTableSession TableChatSessions создание сессии для пользователя
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

// CloseTableSession закрытие сессии для пользователя
func (r *chatRepository) CloseTableSession(sessionID uint) error {
	return r.DB.DB.Model(&storage.TableChatSessions{}).
		Where("id = ?", sessionID).
		Update("active", false).Error
}

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

func (r *chatRepository) UpdateTableSessionActivity(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active": time.Now(),
			"active":      true,
		}).Error
}

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

// Restaurant chat operations

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

func (r *chatRepository) CloseRestaurantSession(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Update("active", false).Error
}

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

func (r *chatRepository) UpdateRestaurantSessionActivity(sessionID uint) error {
	return r.DB.DB.Model(&storage.RestaurantChatSessions{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_active": time.Now(),
			"active":      true,
		}).Error
}

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
