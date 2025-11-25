package service

import (
	"errors"
	orgRepo "savory-ai-server/app/module/organization/repository"
	"savory-ai-server/app/module/question/payload"
	"savory-ai-server/app/module/question/repository"
	"savory-ai-server/app/storage"
)

// questionService реализует бизнес-логику для работы с вопросами.
// Вопросы используются для быстрого способа отправки сообщения в чат бот гостя.
// Каждый вопрос принадлежит организации и может быть привязан к языку для поддержки мультиязычности.
type questionService struct {
	questionRepo repository.QuestionRepository
	languageRepo orgRepo.LanguageRepository
}

// QuestionService определяет контракт для работы с вопросами.
// Все методы работают в контексте организации для обеспечения
// изоляции данных между разными компаниями.
type QuestionService interface {
	GetAll() (*payload.QuestionsResp, error)
	GetByID(id uint) (*payload.QuestionResp, error)
	GetByOrganizationID(id uint) (*payload.QuestionsResp, error)
	GetByOrganizationIDAndLanguage(id uint, languageCode string) (*payload.QuestionsResp, error)
	// GetByFilters возвращает вопросы с фильтрацией по языку и типу чата.
	// Поддерживает комбинации фильтров: только язык, только тип чата, оба или ни одного.
	GetByFilters(organizationID uint, languageCode string, chatType string) (*payload.QuestionsResp, error)
	Create(req *payload.CreateQuestionReq, organizationID uint) (*payload.QuestionResp, error)
	Update(id uint, req *payload.UpdateQuestionReq, organizationID uint) (*payload.QuestionResp, error)
	Delete(id uint, organizationID uint) error
}

// NewQuestionService создаёт новый экземпляр сервиса вопросов.
// Требует репозиторий вопросов и репозиторий языков для работы
// с мультиязычными вопросами.
func NewQuestionService(questionRepo repository.QuestionRepository, languageRepo orgRepo.LanguageRepository) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		languageRepo: languageRepo,
	}
}

// GetAll возвращает все вопросы из системы (без фильтрации по организации).
//
// Бизнес-логика:
// - Используется для административных целей
// - Возвращает вопросы всех организаций с информацией о языке
// - НЕ рекомендуется использовать в пользовательском API (нет изоляции)
func (s *questionService) GetAll() (*payload.QuestionsResp, error) {
	questions, err := s.questionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var questionResps []payload.QuestionResp
	for _, question := range questions {
		questionResps = append(questionResps, mapQuestionToResponse(question))
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

// GetByID возвращает один вопрос по его ID.
//
// Бизнес-логика:
// - Используется для получения детальной информации о вопросе
// - Возвращает вопрос с привязанным языком (если есть)
// - Возвращает ошибку, если вопрос не найден
func (s *questionService) GetByID(id uint) (*payload.QuestionResp, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapQuestionToResponse(question)
	return &resp, nil
}

// GetByOrganizationID возвращает все вопросы конкретной организации.
//
// Бизнес-логика:
// - Основной метод для получения вопросов в пользовательском API
// - Обеспечивает изоляцию данных: организация видит только свои вопросы
// - Возвращает вопросы на всех языках организации
// - Используется, когда язык не указан или нужны все вопросы
func (s *questionService) GetByOrganizationID(id uint) (*payload.QuestionsResp, error) {
	questions, err := s.questionRepo.FindByOrganizationID(id)
	if err != nil {
		return nil, err
	}

	var questionResps []payload.QuestionResp
	for _, question := range questions {
		questionResps = append(questionResps, mapQuestionToResponse(question))
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

// GetByOrganizationIDAndLanguage возвращает вопросы организации,
// отфильтрованные по коду языка.
//
// Бизнес-логика:
// - Используется для показа вопросов гостям на их языке
// - Если languageCode пустой — возвращает все вопросы организации
// - Фильтрует вопросы по коду языка (например: "en", "ru", "es")
// - Возвращает пустой список если нет вопросов на указанном языке
//
// Пример использования:
// - Гость сканирует QR-код и выбирает русский язык
// - Система запрашивает вопросы с languageCode="ru"
// - Гость видит вопросы только на русском языке
func (s *questionService) GetByOrganizationIDAndLanguage(id uint, languageCode string) (*payload.QuestionsResp, error) {
	// Если язык не указан — возвращаем все вопросы организации
	if languageCode == "" {
		return s.GetByOrganizationID(id)
	}

	// Получаем все вопросы организации
	questions, err := s.questionRepo.FindByOrganizationID(id)
	if err != nil {
		return nil, err
	}

	// Фильтруем вопросы по коду языка
	var questionResps []payload.QuestionResp
	for _, question := range questions {
		if question.Language != nil && question.Language.Code == languageCode {
			questionResps = append(questionResps, mapQuestionToResponse(question))
		}
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

// GetByFilters возвращает вопросы с комбинированной фильтрацией по языку и типу чата.
//
// Бизнес-логика:
// - Поддерживает фильтрацию по languageCode (код языка: "en", "ru")
// - Поддерживает фильтрацию по chatType ("reservation" или "menu")
// - Можно комбинировать фильтры или использовать по отдельности
// - Пустые значения игнорируются (не фильтруют)
//
// Примеры:
// - GetByFilters(1, "", "") — все вопросы организации
// - GetByFilters(1, "ru", "") — вопросы на русском языке
// - GetByFilters(1, "", "reservation") — вопросы для чата бронирования
// - GetByFilters(1, "ru", "reservation") — русские вопросы для бронирования
func (s *questionService) GetByFilters(organizationID uint, languageCode string, chatType string) (*payload.QuestionsResp, error) {
	questions, err := s.questionRepo.FindByOrganizationIDLanguageAndChatType(organizationID, languageCode, chatType)
	if err != nil {
		return nil, err
	}

	var questionResps []payload.QuestionResp
	for _, question := range questions {
		questionResps = append(questionResps, mapQuestionToResponse(question))
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

// Create создаёт новый вопрос для организации.
//
// Бизнес-логика:
// - Вопрос автоматически привязывается к организации текущего пользователя
// - Если указан languageCode — вопрос привязывается к языку
// - Если languageCode не указан — вопрос создаётся без языка (универсальный)
// - Возвращает ошибку если указанный язык не существует в системе
//
// Параметры:
// - req.Text: текст вопроса (обязательный)
// - req.LanguageCode: код языка, например "en", "ru" (опциональный)
// - organizationID: ID организации из JWT токена
//
// Пример:
// - Администратор создаёт вопрос "Как вам обслуживание?" с languageCode="ru"
// - Вопрос будет показываться только гостям, выбравшим русский язык
func (s *questionService) Create(req *payload.CreateQuestionReq, organizationID uint) (*payload.QuestionResp, error) {
	// Определяем тип чата (по умолчанию "menu")
	chatType := storage.ChatTypeMenu
	if req.ChatType == "reservation" {
		chatType = storage.ChatTypeReservation
	}

	question := &storage.Question{
		Text:           req.Text,
		OrganizationID: organizationID,
		ChatType:       chatType,
	}

	// Если указан код языка — находим язык и привязываем к вопросу
	if req.LanguageCode != "" {
		language, err := s.languageRepo.FindLanguageByCode(req.LanguageCode)
		if err != nil {
			return nil, errors.New("language not found: " + req.LanguageCode)
		}
		question.LanguageID = &language.ID
	}

	createdQuestion, err := s.questionRepo.Create(question)
	if err != nil {
		return nil, err
	}

	// Перезагружаем вопрос чтобы получить связанный язык
	createdQuestion, err = s.questionRepo.FindByID(createdQuestion.ID)
	if err != nil {
		return nil, err
	}

	resp := mapQuestionToResponse(createdQuestion)
	return &resp, nil
}

// Update обновляет существующий вопрос.
//
// Бизнес-логика:
// - Можно обновить только вопрос своей организации (проверка безопасности)
// - Можно изменить текст вопроса и/или привязку к языку
// - Для смены языка передать новый languageCode
// - Для удаления языка передать пустую строку в languageCode
// - Если поле не передано — оно не изменяется
//
// Параметры:
// - id: ID вопроса для обновления
// - req.Text: новый текст вопроса (опционально)
// - req.LanguageCode: новый код языка или "" для удаления (опционально)
// - organizationID: ID организации из JWT для проверки прав
//
// Ошибки:
// - "question not found" — вопрос с таким ID не существует
// - "question does not belong to your organization" — попытка изменить чужой вопрос
// - "language not found" — указанный язык не существует
func (s *questionService) Update(id uint, req *payload.UpdateQuestionReq, organizationID uint) (*payload.QuestionResp, error) {
	// Находим вопрос по ID
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("question not found")
	}

	// Проверяем что вопрос принадлежит организации пользователя
	if question.OrganizationID != organizationID {
		return nil, errors.New("question does not belong to your organization")
	}

	// Обновляем текст если передан
	if req.Text != "" {
		question.Text = req.Text
	}

	// Обновляем привязку к языку если передана
	if req.LanguageCode != nil {
		if *req.LanguageCode == "" {
			// Пустая строка — удаляем привязку к языку
			question.LanguageID = nil
		} else {
			// Находим язык по коду и привязываем
			language, err := s.languageRepo.FindLanguageByCode(*req.LanguageCode)
			if err != nil {
				return nil, errors.New("language not found: " + *req.LanguageCode)
			}
			question.LanguageID = &language.ID
		}
	}

	// Обновляем тип чата если передан
	if req.ChatType != nil {
		if *req.ChatType == "reservation" {
			question.ChatType = storage.ChatTypeReservation
		} else if *req.ChatType == "menu" {
			question.ChatType = storage.ChatTypeMenu
		}
	}

	// Сохраняем изменения
	updatedQuestion, err := s.questionRepo.Update(question)
	if err != nil {
		return nil, err
	}

	// Перезагружаем для получения связей
	updatedQuestion, err = s.questionRepo.FindByID(updatedQuestion.ID)
	if err != nil {
		return nil, err
	}

	resp := mapQuestionToResponse(updatedQuestion)
	return &resp, nil
}

// Delete удаляет вопрос по ID.
//
// Бизнес-логика:
// - Можно удалить только вопрос своей организации (проверка безопасности)
// - Используется soft delete (GORM устанавливает deleted_at)
// - Удалённые вопросы не возвращаются в списках
//
// Параметры:
// - id: ID вопроса для удаления
// - organizationID: ID организации из JWT для проверки прав
//
// Ошибки:
// - "question not found" — вопрос с таким ID не существует
// - "question does not belong to your organization" — попытка удалить чужой вопрос
func (s *questionService) Delete(id uint, organizationID uint) error {
	// Находим вопрос по ID
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return errors.New("question not found")
	}

	// Проверяем, что вопрос принадлежит организации пользователя
	if question.OrganizationID != organizationID {
		return errors.New("question does not belong to your organization")
	}

	return s.questionRepo.Delete(id)
}

// mapQuestionToResponse преобразует модель Question в DTO для ответа API.
//
// Преобразование:
// - Копирует базовые поля: ID, CreatedAt, Text, ChatType
// - Если к вопросу привязан язык — добавляет информацию о языке
// - Если язык не привязан — поле Language будет nil в JSON (omitempty)
//
// Это позволяет не раскрывать внутреннюю структуру БД клиентам API.
func mapQuestionToResponse(question *storage.Question) payload.QuestionResp {
	resp := payload.QuestionResp{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		Text:      question.Text,
		ChatType:  string(question.ChatType),
	}

	// Добавляем информацию о языке, если она есть
	if question.Language != nil {
		resp.Language = &payload.LanguageResp{
			ID:   question.Language.ID,
			Code: question.Language.Code,
			Name: question.Language.Name,
		}
	}

	return resp
}
