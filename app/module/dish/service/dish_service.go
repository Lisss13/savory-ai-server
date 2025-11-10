package service

import (
	"savory-ai-server/app/module/dish/payload"
	"savory-ai-server/app/module/dish/repository"
	repCategory "savory-ai-server/app/module/menu_category/repository"
	"savory-ai-server/app/storage"
)

type dishService struct {
	dishRepo         repository.DishRepository
	menuCategoryRepo repCategory.MenuCategoryRepository
}

type DishService interface {
	GetDishByMenuCategory() (*payload.DishByCategoryResp, error)
	GetAll() (*payload.DishesResp, error)
	GetByID(id uint) (*payload.DishResp, error)
	GetByOrganizationID(organizationID uint) (*payload.DishesResp, error)
	Create(req *payload.CreateDishReq, organizationID uint) (*payload.DishResp, error)
	Update(id uint, req *payload.UpdateDishReq, organizationID uint) (*payload.DishResp, error)
	Delete(id uint) error
	GetDishOfDay() (*payload.DishResp, error)
	SetDishOfDay(id uint) (*payload.DishResp, error)
}

func NewDishService(dishRepo repository.DishRepository, menuCategoryRepo repCategory.MenuCategoryRepository) DishService {
	return &dishService{
		dishRepo:         dishRepo,
		menuCategoryRepo: menuCategoryRepo,
	}
}

func (s *dishService) GetDishByMenuCategory() (*payload.DishByCategoryResp, error) {
	dishes, err := s.dishRepo.FindAll()
	if err != nil {
		return nil, err
	}
	category, err := s.menuCategoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var dishResps []payload.DishCategoryResp

	for _, cat := range category {
		var dishesInCategory []payload.DishResp
		for _, dish := range dishes {
			if dish.MenuCategoryID == cat.ID {
				dishesInCategory = append(dishesInCategory, mapDishToResponse(dish))
			}
		}
		dishResps = append(dishResps, payload.DishCategoryResp{
			Category: payload.MenuCategoryResp{
				ID:   cat.ID,
				Name: cat.Name,
			},
			Dishes: dishesInCategory,
		})
	}

	return &payload.DishByCategoryResp{
		Dishes: dishResps,
	}, nil
}

func (s *dishService) GetAll() (*payload.DishesResp, error) {
	dishes, err := s.dishRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var dishResps []payload.DishResp
	for _, dish := range dishes {
		dishResps = append(dishResps, mapDishToResponse(dish))
	}

	return &payload.DishesResp{
		Dishes: dishResps,
	}, nil
}

func (s *dishService) GetByOrganizationID(organizationID uint) (*payload.DishesResp, error) {
	dishes, err := s.dishRepo.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var dishResps []payload.DishResp
	for _, dish := range dishes {
		dishResps = append(dishResps, mapDishToResponse(dish))
	}

	return &payload.DishesResp{
		Dishes: dishResps,
	}, nil
}

func (s *dishService) GetByID(id uint) (*payload.DishResp, error) {
	dish, err := s.dishRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(dish)
	return &resp, nil
}

func (s *dishService) Create(req *payload.CreateDishReq, organizationID uint) (*payload.DishResp, error) {
	// Create ingredients
	var ingredients []*storage.Ingredient
	for _, ingredientReq := range req.Ingredients {
		ingredients = append(ingredients, &storage.Ingredient{
			Name:     ingredientReq.Name,
			Quantity: ingredientReq.Quantity,
		})
	}

	// Create a dish
	dish := &storage.Dish{
		OrganizationID: organizationID,
		MenuCategoryID: req.MenuCategoryID,
		Name:           req.Name,
		Price:          req.Price,
		Description:    req.Description,
		Image:          req.Image,
		Ingredients:    ingredients,
	}

	createdDish, err := s.dishRepo.Create(dish)
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(createdDish)
	return &resp, nil
}

func (s *dishService) Update(id uint, req *payload.UpdateDishReq, organizationID uint) (*payload.DishResp, error) {
	// Check if dish exists
	existingDish, err := s.dishRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Create ingredients
	var ingredients []*storage.Ingredient
	for _, ingredientReq := range req.Ingredients {
		ingredients = append(ingredients, &storage.Ingredient{
			DishID:   id,
			Name:     ingredientReq.Name,
			Quantity: ingredientReq.Quantity,
		})
	}

	// Update dish
	existingDish.OrganizationID = organizationID
	existingDish.MenuCategoryID = req.MenuCategoryID
	existingDish.Name = req.Name
	existingDish.Price = req.Price
	existingDish.Description = req.Description
	existingDish.Image = req.Image
	existingDish.Ingredients = ingredients

	updatedDish, err := s.dishRepo.Update(existingDish)
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(updatedDish)
	return &resp, nil
}

func (s *dishService) Delete(id uint) error {
	return s.dishRepo.Delete(id)
}

func (s *dishService) GetDishOfDay() (*payload.DishResp, error) {
	dish, err := s.dishRepo.FindDishOfDay()
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(dish)
	return &resp, nil
}

func (s *dishService) SetDishOfDay(id uint) (*payload.DishResp, error) {
	dish, err := s.dishRepo.SetDishOfDay(id)
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(dish)
	return &resp, nil
}

// Helper function to map a dish to a response
func mapDishToResponse(dish *storage.Dish) payload.DishResp {
	// Map ingredients
	var ingredientResps []payload.IngredientResp
	for _, ingredient := range dish.Ingredients {
		ingredientResps = append(ingredientResps, payload.IngredientResp{
			ID:       ingredient.ID,
			Name:     ingredient.Name,
			Quantity: ingredient.Quantity,
		})
	}

	// Map menu category
	menuCategoryResp := payload.MenuCategoryResp{
		ID:   dish.MenuCategory.ID,
		Name: dish.MenuCategory.Name,
	}

	// Map organization
	organizationResp := payload.OrganizationResp{
		ID:    dish.Organization.ID,
		Name:  dish.Organization.Name,
		Phone: dish.Organization.Phone,
	}

	// Map dish
	return payload.DishResp{
		ID:           dish.ID,
		CreatedAt:    dish.CreatedAt,
		Organization: organizationResp,
		MenuCategory: menuCategoryResp,
		Name:         dish.Name,
		Price:        dish.Price,
		Description:  dish.Description,
		Image:        dish.Image,
		Ingredients:  ingredientResps,
	}
}
