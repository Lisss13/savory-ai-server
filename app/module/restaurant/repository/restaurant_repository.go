package repository

import (
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type restaurantRepository struct {
	DB *database.Database
}

type RestaurantRepository interface {
	FindAll() (restaurants []*storage.Restaurant, err error)
	FindByID(id uint) (restaurant *storage.Restaurant, err error)
	FindByOrganizationID(organizationID uint) (restaurants []*storage.Restaurant, err error)
	Create(restaurant *storage.Restaurant) (res *storage.Restaurant, err error)
	Update(restaurant *storage.Restaurant) (res *storage.Restaurant, err error)
	Delete(id uint) error
}

func NewRestaurantRepository(db *database.Database) RestaurantRepository {
	return &restaurantRepository{
		DB: db,
	}
}

func (r *restaurantRepository) FindAll() (restaurants []*storage.Restaurant, err error) {
	if err := r.DB.DB.Preload("Organization").Preload("WorkingHours").Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *restaurantRepository) FindByID(id uint) (restaurant *storage.Restaurant, err error) {
	err = r.DB.DB.
		Preload("Organization").
		Preload("WorkingHours").
		First(&restaurant, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (r *restaurantRepository) FindByOrganizationID(organizationID uint) (restaurants []*storage.Restaurant, err error) {
	if err := r.DB.DB.Preload("Organization").Preload("WorkingHours").Where("organization_id = ?", organizationID).Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *restaurantRepository) Create(restaurant *storage.Restaurant) (res *storage.Restaurant, err error) {
	if err := r.DB.DB.Create(&restaurant).Error; err != nil {
		return nil, err
	}

	// Reload the restaurant with all associations
	return r.FindByID(restaurant.ID)
}

func (r *restaurantRepository) Update(restaurant *storage.Restaurant) (res *storage.Restaurant, err error) {
	// First, update the restaurant itself
	if err := r.DB.DB.Model(&restaurant).Updates(map[string]interface{}{
		"organization_id": restaurant.OrganizationID,
		"name":            restaurant.Name,
		"address":         restaurant.Address,
		"phone":           restaurant.Phone,
		"website":         restaurant.Website,
		"description":     restaurant.Description,
		"image_url":       restaurant.ImageURL,
	}).Error; err != nil {
		return nil, err
	}

	// If there are working hours, handle them
	//if len(restaurant.WorkingHours) > 0 {
	//	// Delete existing working hours
	//	if err := r.DB.DB.Where("restaurant_id = ?", restaurant.ID).Delete(&storage.WorkingHour{}).Error; err != nil {
	//		return nil, err
	//	}
	//
	//	// Set the restaurant ID for each working hour
	//	for _, workingHour := range restaurant.WorkingHours {
	//		workingHour.RestaurantID = restaurant.ID
	//	}
	//
	//	// Create new working hours
	//	if err := r.DB.DB.Create(&restaurant.WorkingHours).Error; err != nil {
	//		return nil, err
	//	}
	//}

	// Reload the restaurant with all associations
	return r.FindByID(restaurant.ID)
}

func (r *restaurantRepository) Delete(id uint) error {
	// First delete all working hours
	if err := r.DB.DB.Where("restaurant_id = ?", id).Delete(&storage.WorkingHour{}).Error; err != nil {
		return err
	}

	// Then delete the restaurant
	return r.DB.DB.Delete(&storage.Restaurant{}, id).Error
}
