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

// DishService определяет интерфейс сервиса блюд.
type DishService interface {
	GetDishByMenuCategory(restaurantID uint) (*payload.DishByCategoryResp, error)
	GetAll() (*payload.DishesResp, error)
	GetByID(id uint) (*payload.DishResp, error)
	GetByRestaurantID(restaurantID uint) (*payload.DishesResp, error)
	Create(req *payload.CreateDishReq) (*payload.DishResp, error)
	Update(id uint, req *payload.UpdateDishReq) (*payload.DishResp, error)
	Delete(id uint) error
	GetDishOfDay(restaurantID uint) (*payload.DishResp, error)
	SetDishOfDay(id uint) (*payload.DishResp, error)
}

func NewDishService(dishRepo repository.DishRepository, menuCategoryRepo repCategory.MenuCategoryRepository) DishService {
	return &dishService{
		dishRepo:         dishRepo,
		menuCategoryRepo: menuCategoryRepo,
	}
}

// GetDishByMenuCategory возвращает блюда, сгруппированные по категориям для ресторана.
func (s *dishService) GetDishByMenuCategory(restaurantID uint) (*payload.DishByCategoryResp, error) {
	dishes, err := s.dishRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}
	categories, err := s.menuCategoryRepo.FindByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var dishResps []payload.DishCategoryResp

	for _, cat := range categories {
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

// GetByRestaurantID возвращает все блюда для указанного ресторана.
func (s *dishService) GetByRestaurantID(restaurantID uint) (*payload.DishesResp, error) {
	dishes, err := s.dishRepo.FindByRestaurantID(restaurantID)
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

// Create создаёт новое блюдо.
func (s *dishService) Create(req *payload.CreateDishReq) (*payload.DishResp, error) {
	// Create ingredients
	var ingredients []*storage.Ingredient
	for _, ingredientReq := range req.Ingredients {
		ingredients = append(ingredients, &storage.Ingredient{
			Name:     ingredientReq.Name,
			Quantity: ingredientReq.Quantity,
		})
	}

	// Create allergens
	var allergens []*storage.Allergen
	for _, allergenReq := range req.Allergens {
		allergens = append(allergens, &storage.Allergen{
			Name:        allergenReq.Name,
			Description: allergenReq.Description,
		})
	}

	// Create a dish
	dish := &storage.Dish{
		RestaurantID:   req.RestaurantID,
		MenuCategoryID: req.MenuCategoryID,
		Name:           req.Name,
		Price:          req.Price,
		Description:    req.Description,
		Image:          req.Image,
		Ingredients:    ingredients,
		Allergens:      allergens,
	}

	createdDish, err := s.dishRepo.Create(dish)
	if err != nil {
		return nil, err
	}

	resp := mapDishToResponse(createdDish)
	return &resp, nil
}

// Update обновляет существующее блюдо.
func (s *dishService) Update(id uint, req *payload.UpdateDishReq) (*payload.DishResp, error) {
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

	// Create allergens
	var allergens []*storage.Allergen
	for _, allergenReq := range req.Allergens {
		allergens = append(allergens, &storage.Allergen{
			DishID:      id,
			Name:        allergenReq.Name,
			Description: allergenReq.Description,
		})
	}

	// Update dish
	existingDish.RestaurantID = req.RestaurantID
	existingDish.MenuCategoryID = req.MenuCategoryID
	existingDish.Name = req.Name
	existingDish.Price = req.Price
	existingDish.Description = req.Description
	existingDish.Image = req.Image
	existingDish.Ingredients = ingredients
	existingDish.Allergens = allergens

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

// GetDishOfDay возвращает блюдо дня для указанного ресторана.
func (s *dishService) GetDishOfDay(restaurantID uint) (*payload.DishResp, error) {
	dish, err := s.dishRepo.FindDishOfDay(restaurantID)
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

// mapDishToResponse преобразует модель блюда в ответ API.
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

	// Map allergens
	var allergenResps []payload.AllergenResp
	for _, allergen := range dish.Allergens {
		allergenResps = append(allergenResps, payload.AllergenResp{
			ID:          allergen.ID,
			Name:        allergen.Name,
			Description: allergen.Description,
		})
	}

	// Map menu category
	menuCategoryResp := payload.MenuCategoryResp{
		ID:   dish.MenuCategory.ID,
		Name: dish.MenuCategory.Name,
	}

	// Map restaurant
	restaurantResp := payload.RestaurantResp{
		ID:   dish.Restaurant.ID,
		Name: dish.Restaurant.Name,
	}

	// Map dish
	return payload.DishResp{
		ID:           dish.ID,
		CreatedAt:    dish.CreatedAt,
		Restaurant:   restaurantResp,
		MenuCategory: menuCategoryResp,
		Name:         dish.Name,
		Price:        dish.Price,
		Description:  dish.Description,
		Image:        dish.Image,
		Ingredients:  ingredientResps,
		Allergens:    allergenResps,
	}
}
