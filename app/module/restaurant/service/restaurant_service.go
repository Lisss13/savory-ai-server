package service

import (
	"savory-ai-server/app/module/restaurant/payload"
	"savory-ai-server/app/module/restaurant/repository"
	"savory-ai-server/app/storage"
)

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
}

type RestaurantService interface {
	GetAll() (*payload.RestaurantsResp, error)
	GetByID(id uint) (*payload.RestaurantResp, error)
	GetByOrganizationID(organizationID uint) (*payload.RestaurantsResp, error)
	Create(req *payload.CreateRestaurantReq) (*payload.RestaurantResp, error)
	Update(id uint, req *payload.UpdateRestaurantReq) (*payload.RestaurantResp, error)
	Delete(id uint) error
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *restaurantService) GetAll() (*payload.RestaurantsResp, error) {
	restaurants, err := s.restaurantRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var restaurantResps []payload.RestaurantResp
	for _, restaurant := range restaurants {
		restaurantResps = append(restaurantResps, mapRestaurantToResponse(restaurant))
	}

	return &payload.RestaurantsResp{
		Restaurants: restaurantResps,
	}, nil
}

func (s *restaurantService) GetByID(id uint) (*payload.RestaurantResp, error) {
	restaurant, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := mapRestaurantToResponse(restaurant)
	return &resp, nil
}

func (s *restaurantService) GetByOrganizationID(organizationID uint) (*payload.RestaurantsResp, error) {
	restaurants, err := s.restaurantRepo.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var restaurantResps []payload.RestaurantResp
	for _, restaurant := range restaurants {
		restaurantResps = append(restaurantResps, mapRestaurantToResponse(restaurant))
	}

	return &payload.RestaurantsResp{
		Restaurants: restaurantResps,
	}, nil
}

func (s *restaurantService) Create(req *payload.CreateRestaurantReq) (*payload.RestaurantResp, error) {
	// Create working hours
	var workingHours []*storage.WorkingHour
	for _, workingHourReq := range req.WorkingHours {
		workingHours = append(workingHours, &storage.WorkingHour{
			DayOfWeek: workingHourReq.DayOfWeek,
			OpenTime:  workingHourReq.OpenTime,
			CloseTime: workingHourReq.CloseTime,
		})
	}

	// Create restaurant
	restaurant := &storage.Restaurant{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Address:        req.Address,
		Phone:          req.Phone,
		Website:        req.Website,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		WorkingHours:   workingHours,
	}

	createdRestaurant, err := s.restaurantRepo.Create(restaurant)
	if err != nil {
		return nil, err
	}

	resp := mapRestaurantToResponse(createdRestaurant)
	return &resp, nil
}

func (s *restaurantService) Update(id uint, req *payload.UpdateRestaurantReq) (*payload.RestaurantResp, error) {
	// Check if restaurant exists
	existingRestaurant, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update restaurant
	existingRestaurant.OrganizationID = req.OrganizationID
	existingRestaurant.Name = req.Name
	existingRestaurant.Address = req.Address
	existingRestaurant.Phone = req.Phone
	existingRestaurant.Website = req.Website
	existingRestaurant.Description = req.Description
	existingRestaurant.ImageURL = req.ImageURL

	updatedRestaurant, err := s.restaurantRepo.Update(existingRestaurant)
	if err != nil {
		return nil, err
	}

	resp := mapRestaurantToResponse(updatedRestaurant)
	return &resp, nil
}

func (s *restaurantService) Delete(id uint) error {
	return s.restaurantRepo.Delete(id)
}

// Helper function to map a restaurant to a response
func mapRestaurantToResponse(restaurant *storage.Restaurant) payload.RestaurantResp {
	// Map working hours
	var workingHourResps []payload.WorkingHourResp
	for _, workingHour := range restaurant.WorkingHours {
		workingHourResps = append(workingHourResps, payload.WorkingHourResp{
			ID:        workingHour.ID,
			DayOfWeek: workingHour.DayOfWeek,
			OpenTime:  workingHour.OpenTime,
			CloseTime: workingHour.CloseTime,
		})
	}

	// Map user
	organizationResp := payload.OrganizationResp{
		ID:    restaurant.Organization.ID,
		Name:  restaurant.Organization.Name,
		Phone: restaurant.Organization.Phone,
	}

	// Map restaurant
	return payload.RestaurantResp{
		ID:           restaurant.ID,
		CreatedAt:    restaurant.CreatedAt,
		Organization: organizationResp,
		Name:         restaurant.Name,
		Address:      restaurant.Address,
		Phone:        restaurant.Phone,
		Website:      restaurant.Website,
		Description:  restaurant.Description,
		ImageURL:     restaurant.ImageURL,
		WorkingHours: workingHourResps,
	}
}
