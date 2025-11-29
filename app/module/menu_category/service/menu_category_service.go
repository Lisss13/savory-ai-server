package service

import (
	"savory-ai-server/app/module/menu_category/payload"
	"savory-ai-server/app/module/menu_category/repository"
	"savory-ai-server/app/storage"
)

type menuCategoryService struct {
	menuCategoryRepo repository.MenuCategoryRepository
}

// MenuCategoryService определяет интерфейс сервиса категорий меню.
type MenuCategoryService interface {
	// =====================================================
	// CRUD Operations
	// =====================================================
	GetAll() (*payload.MenuCategoriesResp, error)
	GetByID(id uint) (*payload.MenuCategoryResp, error)
	GetByRestaurantID(restaurantID uint) (*payload.MenuCategoriesResp, error)
	Create(req *payload.CreateMenuCategoryReq) (*payload.MenuCategoryResp, error)
	Update(id uint, req *payload.UpdateMenuCategoryReq) (*payload.MenuCategoryResp, error)
	Delete(id uint) error

	// =====================================================
	// Sort Order Operations
	// =====================================================
	// UpdateSortOrder обновляет порядок сортировки для одной категории.
	UpdateSortOrder(id uint, sortOrder int) error
	// UpdateCategoriesSortOrder массово обновляет порядок сортировки категорий.
	UpdateCategoriesSortOrder(req *payload.UpdateCategoriesSortOrderReq) error
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
			ID:           category.ID,
			CreatedAt:    category.CreatedAt,
			Name:         category.Name,
			RestaurantID: category.RestaurantID,
			SortOrder:    category.SortOrder,
		})
	}

	return &payload.MenuCategoriesResp{
		Categories: categoryResps,
	}, nil
}

// GetByRestaurantID возвращает все категории меню для указанного ресторана.
// Категории отсортированы по полю sort_order.
func (s *menuCategoryService) GetByRestaurantID(restaurantID uint) (*payload.MenuCategoriesResp, error) {
	categories, err := s.menuCategoryRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var categoryResps []payload.MenuCategoryResp
	for _, category := range categories {
		categoryResps = append(categoryResps, payload.MenuCategoryResp{
			ID:           category.ID,
			CreatedAt:    category.CreatedAt,
			Name:         category.Name,
			RestaurantID: category.RestaurantID,
			SortOrder:    category.SortOrder,
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
		ID:           category.ID,
		CreatedAt:    category.CreatedAt,
		Name:         category.Name,
		RestaurantID: category.RestaurantID,
		SortOrder:    category.SortOrder,
	}, nil
}

// Create создаёт новую категорию меню для ресторана.
// Если sort_order не указан, автоматически устанавливается следующий порядковый номер.
func (s *menuCategoryService) Create(req *payload.CreateMenuCategoryReq) (*payload.MenuCategoryResp, error) {
	// Определяем sort_order
	var sortOrder int
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	} else {
		// Получаем максимальный sort_order и добавляем 1
		maxSortOrder, err := s.menuCategoryRepo.GetMaxSortOrder(req.RestaurantID)
		if err != nil {
			return nil, err
		}
		sortOrder = maxSortOrder + 1
	}

	category := &storage.MenuCategory{
		Name:         req.Name,
		RestaurantID: req.RestaurantID,
		SortOrder:    sortOrder,
	}

	createdCategory, err := s.menuCategoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return &payload.MenuCategoryResp{
		ID:           createdCategory.ID,
		CreatedAt:    createdCategory.CreatedAt,
		Name:         createdCategory.Name,
		RestaurantID: createdCategory.RestaurantID,
		SortOrder:    createdCategory.SortOrder,
	}, nil
}

// Update обновляет категорию меню.
func (s *menuCategoryService) Update(id uint, req *payload.UpdateMenuCategoryReq) (*payload.MenuCategoryResp, error) {
	category, err := s.menuCategoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем только переданные поля
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	updatedCategory, err := s.menuCategoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return &payload.MenuCategoryResp{
		ID:           updatedCategory.ID,
		CreatedAt:    updatedCategory.CreatedAt,
		Name:         updatedCategory.Name,
		RestaurantID: updatedCategory.RestaurantID,
		SortOrder:    updatedCategory.SortOrder,
	}, nil
}

func (s *menuCategoryService) Delete(id uint) error {
	return s.menuCategoryRepo.Delete(id)
}

// UpdateSortOrder обновляет порядок сортировки для одной категории.
func (s *menuCategoryService) UpdateSortOrder(id uint, sortOrder int) error {
	return s.menuCategoryRepo.UpdateSortOrder(id, sortOrder)
}

// UpdateCategoriesSortOrder массово обновляет порядок сортировки категорий.
// Позволяет пользователю изменить порядок отображения нескольких категорий за один запрос.
func (s *menuCategoryService) UpdateCategoriesSortOrder(req *payload.UpdateCategoriesSortOrderReq) error {
	updates := make(map[uint]int)
	for _, item := range req.Categories {
		updates[item.ID] = item.SortOrder
	}
	return s.menuCategoryRepo.UpdateSortOrderBatch(updates)
}
