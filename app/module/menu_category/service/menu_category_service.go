package service

import (
	"savory-ai-server/app/module/menu_category/payload"
	"savory-ai-server/app/module/menu_category/repository"
	"savory-ai-server/app/storage"
)

type menuCategoryService struct {
	menuCategoryRepo repository.MenuCategoryRepository
}

type MenuCategoryService interface {
	GetAll() (*payload.MenuCategoriesResp, error)
	GetByID(id uint) (*payload.MenuCategoryResp, error)
	GetByOrganizationID(organizationID uint) (*payload.MenuCategoriesResp, error)
	Create(req *payload.CreateMenuCategoryReq, organizationID uint) (*payload.MenuCategoryResp, error)
	Delete(id uint) error
}

func NewMenuCategoryService(menuCategoryRepo repository.MenuCategoryRepository) MenuCategoryService {
	return &menuCategoryService{
		menuCategoryRepo: menuCategoryRepo,
	}
}

func (s *menuCategoryService) GetAll() (*payload.MenuCategoriesResp, error) {
	categories, err := s.menuCategoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var categoryResps []payload.MenuCategoryResp
	for _, category := range categories {
		categoryResps = append(categoryResps, payload.MenuCategoryResp{
			ID:        category.ID,
			CreatedAt: category.CreatedAt,
			Name:      category.Name,
		})
	}

	return &payload.MenuCategoriesResp{
		Categories: categoryResps,
	}, nil
}

func (s *menuCategoryService) GetByOrganizationID(organizationID uint) (*payload.MenuCategoriesResp, error) {
	categories, err := s.menuCategoryRepo.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var categoryResps []payload.MenuCategoryResp
	for _, category := range categories {
		categoryResps = append(categoryResps, payload.MenuCategoryResp{
			ID:        category.ID,
			CreatedAt: category.CreatedAt,
			Name:      category.Name,
		})
	}

	return &payload.MenuCategoriesResp{
		Categories: categoryResps,
	}, nil
}

func (s *menuCategoryService) GetByID(id uint) (*payload.MenuCategoryResp, error) {
	category, err := s.menuCategoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &payload.MenuCategoryResp{
		ID:        category.ID,
		CreatedAt: category.CreatedAt,
		Name:      category.Name,
	}, nil
}

func (s *menuCategoryService) Create(req *payload.CreateMenuCategoryReq, organizationID uint) (*payload.MenuCategoryResp, error) {
	category := &storage.MenuCategory{
		Name:           req.Name,
		OrganizationID: organizationID,
	}

	createdCategory, err := s.menuCategoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return &payload.MenuCategoryResp{
		ID:        createdCategory.ID,
		CreatedAt: createdCategory.CreatedAt,
		Name:      createdCategory.Name,
	}, nil
}

func (s *menuCategoryService) Delete(id uint) error {
	return s.menuCategoryRepo.Delete(id)
}
