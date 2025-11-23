package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type organizationRepository struct {
	DB *database.Database
}

type OrganizationRepository interface {
	FindOrganizationByID(id uint) (*storage.Organization, error)
	FindAllOrganizations() ([]*storage.Organization, error)
	CreateOrganization(org *storage.Organization) (*storage.Organization, error)
	UpdateOrganization(org *storage.Organization) (*storage.Organization, error)
	AddUserToOrganization(orgID, userID uint) error
	RemoveUserFromOrganization(orgID, userID uint) error
}

func NewOrganizationRepository(db *database.Database) OrganizationRepository {
	return &organizationRepository{
		DB: db,
	}
}

func (or *organizationRepository) FindOrganizationByID(id uint) (*storage.Organization, error) {
	var org *storage.Organization
	err := or.DB.DB.
		Preload("Admin").
		Preload("Users").
		Preload("Languages").
		First(&org, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (or *organizationRepository) FindAllOrganizations() ([]*storage.Organization, error) {
	var orgs []*storage.Organization
	err := or.DB.DB.
		Preload("Admin").
		Find(&orgs).
		Error

	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (or *organizationRepository) CreateOrganization(org *storage.Organization) (*storage.Organization, error) {
	if err := or.DB.DB.Create(&org).Error; err != nil {
		return nil, err
	}

	return org, nil
}

func (or *organizationRepository) UpdateOrganization(org *storage.Organization) (*storage.Organization, error) {
	if err := or.DB.DB.
		Model(&org).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"name":  org.Name,
			"phone": org.Phone,
		}).Error; err != nil {
		return nil, err
	}

	return org, nil
}

func (or *organizationRepository) FindOrganizationsByUserID(userID uint) ([]*storage.Organization, error) {
	var orgs []*storage.Organization
	err := or.DB.DB.
		Model(&storage.Organization{}).
		Joins("JOIN organization_users ou ON ou.organization_id = organizations.id").
		Where("ou.user_id = ?", userID).
		Preload("Admin").
		Preload("Users").
		Find(&orgs).Error

	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (or *organizationRepository) AddUserToOrganization(orgID, userID uint) error {
	return or.DB.DB.Exec("INSERT INTO organization_users (organization_id, user_id) VALUES (?, ?)", orgID, userID).Error
}

func (or *organizationRepository) RemoveUserFromOrganization(orgID, userID uint) error {
	return or.DB.DB.Exec("DELETE FROM organization_users WHERE organization_id = ? AND user_id = ?", orgID, userID).Error
}
