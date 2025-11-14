package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"savory-ai-server/app/module/auth/payload"
	organizationRepo "savory-ai-server/app/module/organization/repository"
	organizationService "savory-ai-server/app/module/organization/service"
	qrCodeService "savory-ai-server/app/module/qr_code/service"
	user_repo "savory-ai-server/app/module/user/repository"
	"savory-ai-server/app/storage"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/jwt"
)

// AuthService
type authService struct {
	userRepo            user_repo.UserRepository
	config              *config.Config
	qrService           qrCodeService.QRCodeService
	organizationRepo    organizationRepo.OrganizationRepository
	organizationService organizationService.OrganizationService
}

type AuthService interface {
	Login(req payload.LoginRequest) (res payload.LoginResponse, err error)
	Register(req payload.RegisterRequest) (res payload.RegisterResponse, err error)
}

// AuthService
func NewAuthService(
	userRepo user_repo.UserRepository,
	cfg *config.Config,
	qrService qrCodeService.QRCodeService,
	organizationRepo organizationRepo.OrganizationRepository,
	organizationService organizationService.OrganizationService,
) AuthService {
	return &authService{
		userRepo:            userRepo,
		config:              cfg,
		qrService:           qrService,
		organizationService: organizationService,
		organizationRepo:    organizationRepo,
	}
}

func (as *authService) Login(req payload.LoginRequest) (res payload.LoginResponse, err error) {
	// check user by email
	user, err := as.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return
	}

	if user == nil {
		err = errors.New("user not found")
		return
	}

	// check password
	if !user.ComparePassword(req.Password) {
		err = errors.New("password not match")
		return
	}

	company, err := as.organizationRepo.FindOrganizationByID(user.ID)
	if err != nil {
		return
	}

	// do create token
	token, exp, err := jwt.
		NewJWT(as.config.Middleware.Jwt.Secret, as.config.Middleware.Jwt.Expiration).
		GenerateToken(jwt.JWTData{
			ID:        user.ID,
			Email:     user.Email,
			CompanyID: company.ID,
		})
	if err != nil {
		return
	}

	res.Token = token
	res.Type = "Bearer"
	res.ExpiresAt = exp.Unix()
	res.User = payload.RegisterResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Company: user.Company,
		Phone:   user.Phone,
	}

	res.Organization = payload.UserOrganizationResponse{
		ID:      company.ID,
		Name:    company.Name,
		Company: company.Name,
		Phone:   company.Phone,
		AdminID: company.AdminID,
	}
	fmt.Println("res.Organization", res.Organization)

	return
}

func (as *authService) Register(req payload.RegisterRequest) (res payload.RegisterResponse, err error) {
	// check user by email
	user, err := as.userRepo.FindUserByEmail(req.Email)

	if user != nil {
		err = errors.New("email already exists")
		return
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	// do create user
	newUser := &storage.User{
		Company:  req.CompanyName,
		Email:    req.Email,
		Phone:    req.Phone,
		Name:     req.Name,
		Password: string(bcryptPassword),
	}

	user, err = as.userRepo.CreateUser(newUser)
	if err != nil {
		return
	}

	// Создание компании и связывание с пользователем
	if err = as.organizationService.CreateOrganization(user.ID, req.CompanyName, req.Phone); err != nil {
		return
	}

	// generate qr code for user
	if _, err = as.qrService.GenerateUserQRCode(user.ID); err != nil {
		return
	}

	res.ID = user.ID
	res.Email = user.Email
	res.Name = user.Name
	res.Company = user.Company
	res.Phone = user.Phone

	return
}
