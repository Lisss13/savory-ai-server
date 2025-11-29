package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type questionRepository struct {
	DB *database.Database
}

// QuestionRepository определяет контракт для работы с вопросами в БД.
type QuestionRepository interface {
	FindAll() (questions []*storage.Question, err error)
	FindByID(id uint) (question *storage.Question, err error)
	FindByOrganizationID(organizationID uint) (questions []*storage.Question, err error)
	FindByOrganizationIDAndLanguage(organizationID uint, languageID *uint) (questions []*storage.Question, err error)
	// FindByOrganizationIDAndChatType возвращает вопросы организации,
	// отфильтрованные по типу чата (reservation/menu).
	FindByOrganizationIDAndChatType(organizationID uint, chatType string) (questions []*storage.Question, err error)
	// FindByOrganizationIDLanguageAndChatType возвращает вопросы с фильтрацией
	// по языку и типу чата одновременно.
	FindByOrganizationIDLanguageAndChatType(organizationID uint, languageCode string, chatType string) (questions []*storage.Question, err error)
	Create(question *storage.Question) (res *storage.Question, err error)
	Update(question *storage.Question) (res *storage.Question, err error)
	Delete(id uint) error
	// GetMaxDisplayOrder возвращает максимальный display_order для организации.
	// Используется при создании нового вопроса для автоинкремента порядка.
	GetMaxDisplayOrder(organizationID uint) (int, error)
	// UpdateDisplayOrder обновляет порядок отображения для вопроса.
	UpdateDisplayOrder(questionID uint, displayOrder int) error
}

// NewQuestionRepository создаёт новый экземпляр репозитория вопросов.
func NewQuestionRepository(db *database.Database) QuestionRepository {
	return &questionRepository{
		DB: db,
	}
}

func (r *questionRepository) FindAll() (questions []*storage.Question, err error) {
	if err := r.DB.DB.Preload("Organization").Preload("Language").Order("display_order ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) FindByOrganizationID(organizationID uint) (questions []*storage.Question, err error) {
	if err := r.DB.DB.Preload("Organization").Preload("Language").Where("organization_id = ?", organizationID).Order("display_order ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) FindByOrganizationIDAndLanguage(organizationID uint, languageID *uint) (questions []*storage.Question, err error) {
	query := r.DB.DB.Preload("Organization").Preload("Language").Where("organization_id = ?", organizationID)

	if languageID != nil {
		query = query.Where("language_id = ?", *languageID)
	} else {
		query = query.Where("language_id IS NULL")
	}

	if err := query.Order("display_order ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) FindByID(id uint) (question *storage.Question, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload("Language").
		First(&question, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return question, nil
}

func (r *questionRepository) Create(question *storage.Question) (res *storage.Question, err error) {
	if err := r.DB.DB.Create(&question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (r *questionRepository) Update(question *storage.Question) (res *storage.Question, err error) {
	if err := r.DB.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (r *questionRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.Question{}, id).Error
}

// FindByOrganizationIDAndChatType возвращает вопросы организации по типу чата.
// Если chatType пустой — возвращает все вопросы организации.
// Результаты отсортированы по display_order (по возрастанию).
func (r *questionRepository) FindByOrganizationIDAndChatType(organizationID uint, chatType string) (questions []*storage.Question, err error) {
	query := r.DB.DB.Preload("Organization").Preload("Language").Where("organization_id = ?", organizationID)

	if chatType != "" {
		query = query.Where("chat_type = ?", chatType)
	}

	if err := query.Order("display_order ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// FindByOrganizationIDLanguageAndChatType возвращает вопросы с комбинированной фильтрацией.
// Фильтрует по организации, коду языка и типу чата.
// Пустые значения languageCode и chatType игнорируются при фильтрации.
// Результаты отсортированы по display_order (по возрастанию).
func (r *questionRepository) FindByOrganizationIDLanguageAndChatType(organizationID uint, languageCode string, chatType string) (questions []*storage.Question, err error) {
	query := r.DB.DB.Preload("Organization").Preload("Language").Where("organization_id = ?", organizationID)

	if chatType != "" {
		query = query.Where("chat_type = ?", chatType)
	}

	if err := query.Order("display_order ASC").Find(&questions).Error; err != nil {
		return nil, err
	}

	// Если указан код языка — фильтруем в памяти (для поддержки фильтрации по Language.Code)
	if languageCode != "" {
		var filtered []*storage.Question
		for _, q := range questions {
			if q.Language != nil && q.Language.Code == languageCode {
				filtered = append(filtered, q)
			}
		}
		return filtered, nil
	}

	return questions, nil
}

// GetMaxDisplayOrder возвращает максимальный display_order для организации.
// Если вопросов нет, возвращает -1 (чтобы первый вопрос получил 0).
func (r *questionRepository) GetMaxDisplayOrder(organizationID uint) (int, error) {
	var maxOrder *int
	err := r.DB.DB.Model(&storage.Question{}).
		Where("organization_id = ?", organizationID).
		Select("MAX(display_order)").
		Scan(&maxOrder).Error
	if err != nil {
		return 0, err
	}
	if maxOrder == nil {
		return -1, nil
	}
	return *maxOrder, nil
}

// UpdateDisplayOrder обновляет порядок отображения для вопроса.
func (r *questionRepository) UpdateDisplayOrder(questionID uint, displayOrder int) error {
	return r.DB.DB.Model(&storage.Question{}).
		Where("id = ?", questionID).
		Update("display_order", displayOrder).Error
}
