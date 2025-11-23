package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type subscriptionRepository struct {
	DB *database.Database
}

type SubscriptionRepository interface {
	FindAll() (subscriptions []*storage.Subscription, err error)
	FindByID(id uint) (subscription *storage.Subscription, err error)
	FindByOrganizationID(organizationID uint) (subscriptions []*storage.Subscription, err error)
	FindActiveByOrganizationID(organizationID uint) (subscription *storage.Subscription, err error)
	Create(subscription *storage.Subscription) (res *storage.Subscription, err error)
	Update(subscription *storage.Subscription) (res *storage.Subscription, err error)
	Delete(id uint) error
}

func NewSubscriptionRepository(db *database.Database) SubscriptionRepository {
	return &subscriptionRepository{
		DB: db,
	}
}

func (r *subscriptionRepository) FindAll() (subscriptions []*storage.Subscription, err error) {
	if err := r.DB.DB.Preload("Organization").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *subscriptionRepository) FindByID(id uint) (subscription *storage.Subscription, err error) {
	err = r.DB.DB.
		Preload("Organization").
		First(&subscription, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (r *subscriptionRepository) FindByOrganizationID(organizationID uint) (subscriptions []*storage.Subscription, err error) {
	if err := r.DB.DB.Preload("Organization").Where("organization_id = ?", organizationID).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *subscriptionRepository) FindActiveByOrganizationID(organizationID uint) (subscription *storage.Subscription, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Where("organization_id = ? AND is_active = ?", organizationID, true).
		First(&subscription).
		Error

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (r *subscriptionRepository) Create(subscription *storage.Subscription) (res *storage.Subscription, err error) {
	if err := r.DB.DB.Create(&subscription).Error; err != nil {
		return nil, err
	}

	// Reload the subscription with all associations
	return r.FindByID(subscription.ID)
}

func (r *subscriptionRepository) Update(subscription *storage.Subscription) (res *storage.Subscription, err error) {
	if err := r.DB.DB.Model(&subscription).Updates(map[string]interface{}{
		"period":     subscription.Period,
		"start_date": subscription.StartDate,
		"end_date":   subscription.EndDate,
		"is_active":  subscription.IsActive,
	}).Error; err != nil {
		return nil, err
	}

	// Reload the subscription with all associations
	return r.FindByID(subscription.ID)
}

func (r *subscriptionRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.Subscription{}, id).Error
}
