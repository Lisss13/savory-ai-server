package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type languageRepository struct {
	DB *database.Database
}

type LanguageRepository interface {
	FindAllLanguages() ([]*storage.Language, error)
	FindLanguageByID(id uint) (*storage.Language, error)
	FindLanguageByCode(code string) (*storage.Language, error)
	CreateLanguage(language *storage.Language) (*storage.Language, error)
	UpdateLanguage(language *storage.Language) (*storage.Language, error)
	DeleteLanguage(id uint) error
	
	// Organization language operations
	FindLanguagesByOrganizationID(orgID uint) ([]*storage.Language, error)
	AddLanguageToOrganization(orgID, langID uint) error
	RemoveLanguageFromOrganization(orgID, langID uint) error
}

func NewLanguageRepository(db *database.Database) LanguageRepository {
	return &languageRepository{
		DB: db,
	}
}

func (lr *languageRepository) FindAllLanguages() ([]*storage.Language, error) {
	var languages []*storage.Language
	err := lr.DB.DB.
		Find(&languages).
		Error

	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (lr *languageRepository) FindLanguageByID(id uint) (*storage.Language, error) {
	var language *storage.Language
	err := lr.DB.DB.
		First(&language, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return language, nil
}

func (lr *languageRepository) FindLanguageByCode(code string) (*storage.Language, error) {
	var language *storage.Language
	err := lr.DB.DB.
		First(&language, "code = ?", code).
		Error

	if err != nil {
		return nil, err
	}

	return language, nil
}

func (lr *languageRepository) CreateLanguage(language *storage.Language) (*storage.Language, error) {
	if err := lr.DB.DB.Create(&language).Error; err != nil {
		return nil, err
	}

	return language, nil
}

func (lr *languageRepository) UpdateLanguage(language *storage.Language) (*storage.Language, error) {
	if err := lr.DB.DB.
		Model(&language).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"code":        language.Code,
			"name":        language.Name,
			"description": language.Description,
		}).Error; err != nil {
		return nil, err
	}

	return language, nil
}

func (lr *languageRepository) DeleteLanguage(id uint) error {
	return lr.DB.DB.Delete(&storage.Language{}, id).Error
}

func (lr *languageRepository) FindLanguagesByOrganizationID(orgID uint) ([]*storage.Language, error) {
	var languages []*storage.Language
	err := lr.DB.DB.
		Model(&storage.Language{}).
		Joins("JOIN organization_languages ol ON ol.language_id = languages.id").
		Where("ol.organization_id = ?", orgID).
		Find(&languages).Error

	if err != nil {
		return nil, err
	}

	return languages, nil
}

func (lr *languageRepository) AddLanguageToOrganization(orgID, langID uint) error {
	return lr.DB.DB.Exec("INSERT INTO organization_languages (organization_id, language_id) VALUES (?, ?)", orgID, langID).Error
}

func (lr *languageRepository) RemoveLanguageFromOrganization(orgID, langID uint) error {
	return lr.DB.DB.Exec("DELETE FROM organization_languages WHERE organization_id = ? AND language_id = ?", orgID, langID).Error
}