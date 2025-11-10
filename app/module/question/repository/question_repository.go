package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type questionRepository struct {
	DB *database.Database
}

type QuestionRepository interface {
	FindAll() (questions []*storage.Question, err error)
	FindByID(id uint) (question *storage.Question, err error)
	FindByOrganizationID(organizationID uint) (questions []*storage.Question, err error)
	Create(question *storage.Question) (res *storage.Question, err error)
	Delete(id uint) error
}

func NewQuestionRepository(db *database.Database) QuestionRepository {
	return &questionRepository{
		DB: db,
	}
}

func (r *questionRepository) FindAll() (questions []*storage.Question, err error) {
	if err := r.DB.DB.Preload("Organization").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) FindByOrganizationID(organizationID uint) (questions []*storage.Question, err error) {
	if err := r.DB.DB.Preload("Organization").Where("organization_id = ?", organizationID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionRepository) FindByID(id uint) (question *storage.Question, err error) {
	err = r.DB.DB.
		Preload("Organization").
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

func (r *questionRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.Question{}, id).Error
}
