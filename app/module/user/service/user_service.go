package service

import (
	"golang.org/x/crypto/bcrypt"
	"savory-ai-server/app/module/organization/service"
	"savory-ai-server/app/module/user/payload"
	"savory-ai-server/app/module/user/repository"
	"savory-ai-server/app/storage"
)

type userService struct {
	userRepo            repository.UserRepository
	organizationService service.OrganizationService
}

type UserService interface {
	FindUserByID(userID int64) (*payload.UserResp, error)
	CreateUser(req *payload.UserCreateReq, CompanyID uint) (*payload.UserResp, error)
}

func NewUserService(userRepo repository.UserRepository, organizationService service.OrganizationService) UserService {
	return &userService{
		userRepo:            userRepo,
		organizationService: organizationService,
	}
}

func (ur *userService) FindUserByID(userID int64) (*payload.UserResp, error) {
	user, err := ur.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	resp := &payload.UserResp{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		Name:      user.Name,
		Company:   user.Company,
		Phone:     user.Phone,
	}
	return resp, nil
}

func (ur *userService) CreateUser(req *payload.UserCreateReq, companyID uint) (*payload.UserResp, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &storage.User{
		Name:     req.Name,
		Company:  req.Company,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	createdUser, err := ur.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// add user to an organization
	err = ur.organizationService.AddUserToOrganization(companyID, createdUser.ID)
	if err != nil {
		return nil, err
	}

	// Return the user response
	resp := &payload.UserResp{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		Email:     createdUser.Email,
		Name:      createdUser.Name,
		Company:   createdUser.Company,
		Phone:     createdUser.Phone,
	}
	return resp, nil
}
