package service

import (
	"savory-ai-server/app/module/question/payload"
	"savory-ai-server/app/module/question/repository"
	"savory-ai-server/app/storage"
)

type questionService struct {
	questionRepo repository.QuestionRepository
}

type QuestionService interface {
	GetAll() (*payload.QuestionsResp, error)
	GetByID(id uint) (*payload.QuestionResp, error)
	GetByOrganizationID(id uint) (*payload.QuestionsResp, error)
	Create(req *payload.CreateQuestionReq, organizationID uint) (*payload.QuestionResp, error)
	Delete(id uint) error
}

func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
	}
}

func (s *questionService) GetAll() (*payload.QuestionsResp, error) {
	questions, err := s.questionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var questionResps []payload.QuestionResp
	for _, question := range questions {
		questionResps = append(questionResps, payload.QuestionResp{
			ID:        question.ID,
			CreatedAt: question.CreatedAt,
			Text:      question.Text,
		})
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

func (s *questionService) GetByID(id uint) (*payload.QuestionResp, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &payload.QuestionResp{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		Text:      question.Text,
	}, nil
}

func (s *questionService) GetByOrganizationID(id uint) (*payload.QuestionsResp, error) {
	questions, err := s.questionRepo.FindByOrganizationID(id)
	if err != nil {
		return nil, err
	}

	var questionResps []payload.QuestionResp
	for _, question := range questions {
		questionResps = append(questionResps, payload.QuestionResp{
			ID:        question.ID,
			CreatedAt: question.CreatedAt,
			Text:      question.Text,
		})
	}

	return &payload.QuestionsResp{
		Questions: questionResps,
	}, nil
}

func (s *questionService) Create(req *payload.CreateQuestionReq, organizationID uint) (*payload.QuestionResp, error) {
	question := &storage.Question{
		Text:           req.Text,
		OrganizationID: organizationID,
	}

	createdQuestion, err := s.questionRepo.Create(question)
	if err != nil {
		return nil, err
	}

	return &payload.QuestionResp{
		ID:        createdQuestion.ID,
		CreatedAt: createdQuestion.CreatedAt,
		Text:      createdQuestion.Text,
	}, nil
}

func (s *questionService) Delete(id uint) error {
	return s.questionRepo.Delete(id)
}
