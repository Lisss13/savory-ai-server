package service

import (
	qrCodeService "savory-ai-server/app/module/qr_code/service"
	"savory-ai-server/app/module/table/payload"
	"savory-ai-server/app/module/table/repository"
	"savory-ai-server/app/storage"
)

type tableService struct {
	tableRepo repository.TableRepository
	qrService qrCodeService.QRCodeService
}

type TableService interface {
	GetAll() (*payload.TablesResp, error)
	GetByID(id uint) (*payload.TableResp, error)
	GetByRestaurantID(restaurantID uint) (*payload.TablesResp, error)
	Create(req *payload.CreateTableReq) (*payload.TableResp, error)
	Update(id uint, req *payload.UpdateTableReq) (*payload.TableResp, error)
	Delete(id uint) error
}

func NewTableService(tableRepo repository.TableRepository, qrService qrCodeService.QRCodeService) TableService {
	return &tableService{
		tableRepo: tableRepo,
		qrService: qrService,
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

func (s *tableService) Create(req *payload.CreateTableReq) (*payload.TableResp, error) {
	// Create table
	table := &storage.Table{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		GuestCount:   req.GuestCount,
	}

	createdTable, err := s.tableRepo.Create(table)
	if err != nil {
		return nil, err
	}

	if _, err = s.qrService.GenerateTableQRCode(table.RestaurantID, table.ID); err != nil {
		return nil, err
	}

	resp := mapTableToResponse(createdTable)
	return &resp, nil
}

func (s *tableService) Update(id uint, req *payload.UpdateTableReq) (*payload.TableResp, error) {
	// Check if table exists
	existingTable, err := s.tableRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update table
	existingTable.RestaurantID = req.RestaurantID
	existingTable.Name = req.Name
	existingTable.GuestCount = req.GuestCount

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
	}
}
