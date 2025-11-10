package service

import (
	"fmt"
	"savory-ai-server/app/module/table/payload"
	"savory-ai-server/app/module/table/repository"
	"savory-ai-server/app/storage"
)

type tableService struct {
	tableRepo repository.TableRepository
}

type TableService interface {
	GetAll() (*payload.TablesResp, error)
	GetByID(id uint) (*payload.TableResp, error)
	GetByRestaurantID(restaurantID uint) (*payload.TablesResp, error)
	Create(req *payload.CreateTableReq, restaurantID uint) (*payload.TableResp, error)
	Update(id uint, req *payload.UpdateTableReq, restaurantID uint) (*payload.TableResp, error)
	Delete(id uint) error
	GenerateQRCodeURL(restaurantID, tableID uint) string
}

func NewTableService(tableRepo repository.TableRepository) TableService {
	return &tableService{
		tableRepo: tableRepo,
	}
}

func (s *tableService) GetAll() (*payload.TablesResp, error) {
	tables, err := s.tableRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var tableResps []payload.TableResp
	for _, table := range tables {
		tableResps = append(tableResps, mapTableToResponse(table))
	}

	return &payload.TablesResp{
		Tables: tableResps,
	}, nil
}

func (s *tableService) GetByID(id uint) (*payload.TableResp, error) {
	table, err := s.tableRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapTableToResponse(table)
	return &resp, nil
}

func (s *tableService) GetByRestaurantID(restaurantID uint) (*payload.TablesResp, error) {
	tables, err := s.tableRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var tableResps []payload.TableResp
	for _, table := range tables {
		tableResps = append(tableResps, mapTableToResponse(table))
	}

	return &payload.TablesResp{
		Tables: tableResps,
	}, nil
}

func (s *tableService) Create(req *payload.CreateTableReq, restaurantID uint) (*payload.TableResp, error) {
	// Create table
	table := &storage.Table{
		RestaurantID: restaurantID,
		Name:         req.Name,
		GuestCount:   req.GuestCount,
	}

	// Generate QR code URL
	table.QRCodeURL = s.GenerateQRCodeURL(restaurantID, 0) // 0 will be replaced with the actual ID after creation

	createdTable, err := s.tableRepo.Create(table)
	if err != nil {
		return nil, err
	}

	// Update QR code URL with the actual table ID
	createdTable.QRCodeURL = s.GenerateQRCodeURL(restaurantID, createdTable.ID)
	updatedTable, err := s.tableRepo.Update(createdTable)
	if err != nil {
		return nil, err
	}

	resp := mapTableToResponse(updatedTable)
	return &resp, nil
}

func (s *tableService) Update(id uint, req *payload.UpdateTableReq, restaurantID uint) (*payload.TableResp, error) {
	// Check if table exists
	existingTable, err := s.tableRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update table
	existingTable.RestaurantID = restaurantID
	existingTable.Name = req.Name
	existingTable.GuestCount = req.GuestCount
	existingTable.QRCodeURL = s.GenerateQRCodeURL(restaurantID, id)

	updatedTable, err := s.tableRepo.Update(existingTable)
	if err != nil {
		return nil, err
	}

	resp := mapTableToResponse(updatedTable)
	return &resp, nil
}

func (s *tableService) Delete(id uint) error {
	return s.tableRepo.Delete(id)
}

func (s *tableService) GenerateQRCodeURL(restaurantID, tableID uint) string {
	if tableID == 0 {
		// This is a placeholder URL that will be updated after table creation
		return fmt.Sprintf("http://localhost:5000/restaurant/%d/table/placeholder", restaurantID)
	}
	return fmt.Sprintf("http://localhost:5000/restaurant/%d/table/%d", restaurantID, tableID)
}

// Helper function to map a table to a response
func mapTableToResponse(table *storage.Table) payload.TableResp {
	// Map restaurant
	restaurantResp := payload.RestaurantResp{
		ID:   table.Restaurant.ID,
		Name: table.Restaurant.Name,
	}

	// Map table
	return payload.TableResp{
		ID:         table.ID,
		CreatedAt:  table.CreatedAt,
		Restaurant: restaurantResp,
		Name:       table.Name,
		GuestCount: table.GuestCount,
		QRCodeURL:  table.QRCodeURL,
	}
}
