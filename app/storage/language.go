package storage

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Language represents a language that can be associated with an organization
type Language struct {
	gorm.Model
	Code        string `gorm:"column:code;not null;uniqueIndex" json:"code"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	// Organizations that use this language
	Organizations []Organization `gorm:"many2many:organization_languages;" json:"organizations"`
}

// OrganizationLanguage represents the many-to-many relationship between organizations and languages
type OrganizationLanguage struct {
	OrganizationID uint `gorm:"primaryKey"`
	LanguageID     uint `gorm:"primaryKey"`
}

// LanguageSeeder is responsible for seeding default languages into the database
type LanguageSeeder struct {
	DB *gorm.DB
}

// Count returns the number of languages in the database
func (s *LanguageSeeder) Count() (int, error) {
	var cnt int64
	if err := s.DB.Model(&Language{}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return int(cnt), nil
}

// Seed adds default languages to the database
// если они удалены пользователем, а не в ручную, то в базу они не перезапишутся,
// если хочешь перезаписать, надо вручную удалить из базы и перезапустить сидер
func (s *LanguageSeeder) Seed(db *gorm.DB) error {
	languages := []Language{
		{Code: "en", Name: "English", Description: "English Language"},
		{Code: "ru", Name: "Русский", Description: "Russian Language"},
	}
	// если языки уже существуют, не добавлять их снова,
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&languages).Error
}
