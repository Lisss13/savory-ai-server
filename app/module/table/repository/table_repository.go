package repository

import (
	"gorm.io/gorm/clause"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
)

type tableRepository struct {
	DB *database.Database
}

type TableRepository interface {
	FindAll() (tables []*storage.Table, err error)
	FindByID(id uint) (table *storage.Table, err error)
	FindByRestaurantID(restaurantID uint) (tables []*storage.Table, err error)
	Create(table *storage.Table) (res *storage.Table, err error)
	Update(table *storage.Table) (res *storage.Table, err error)
	Delete(id uint) error
}

func NewTableRepository(db *database.Database) TableRepository {
	return &tableRepository{
		DB: db,
	}
}

func (r *tableRepository) FindAll() (tables []*storage.Table, err error) {
	if err := r.DB.DB.Preload("Restaurant").Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *tableRepository) FindByID(id uint) (table *storage.Table, err error) {
	err = r.DB.DB.
		Preload(clause.Associations).
		First(&table, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return table, nil
}

func (r *tableRepository) FindByRestaurantID(restaurantID uint) (tables []*storage.Table, err error) {
	if err := r.DB.DB.Preload("Restaurant").Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *tableRepository) Create(table *storage.Table) (res *storage.Table, err error) {
	if err := r.DB.DB.Create(&table).Error; err != nil {
		return nil, err
	}

	// Reload the table with all associations
	return r.FindByID(table.ID)
}

func (r *tableRepository) Update(table *storage.Table) (res *storage.Table, err error) {
	// Update the table
	if err := r.DB.DB.Model(&table).Updates(map[string]interface{}{
		"restaurant_id": table.RestaurantID,
		"name":          table.Name,
		"guest_count":   table.GuestCount,
		"qr_code_url":   table.QRCodeURL,
	}).Error; err != nil {
		return nil, err
	}

	// Reload the table with all associations
	return r.FindByID(table.ID)
}

func (r *tableRepository) Delete(id uint) error {
	return r.DB.DB.Delete(&storage.Table{}, id).Error
}
