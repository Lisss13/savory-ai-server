package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"savory-ai-server/app/module/auth/payload"
	organizationRepo "savory-ai-server/app/module/organization/repository"
	organizationService "savory-ai-server/app/module/organization/service"
	user_repo "savory-ai-server/app/module/user/repository"
	"savory-ai-server/app/storage"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/jwt"
	"strconv"
	"time"
)

// AuthService
type authService struct {
	userRepo            user_repo.UserRepository
	config              *config.Config
	organizationRepo    organizationRepo.OrganizationRepository
	organizationService organizationService.OrganizationService
}

type AuthService interface {
	Login(req payload.LoginRequest) (res payload.LoginResponse, err error)
	Register(req payload.RegisterRequest) (res payload.RegisterResponse, err error)
	ChangePassword(userID uint, req payload.ChangePasswordRequest) (res payload.ChangePasswordResponse, err error)

	// Password reset methods
	RequestPasswordReset(req payload.RequestPasswordResetRequest) (res payload.RequestPasswordResetResponse, err error)
	VerifyPasswordReset(req payload.VerifyPasswordResetRequest) (res payload.VerifyPasswordResetResponse, err error)
}

// AuthService
func NewAuthService(
	userRepo user_repo.UserRepository,
	cfg *config.Config,
	organizationRepo organizationRepo.OrganizationRepository,
	organizationService organizationService.OrganizationService,
) AuthService {
	return &authService{
		userRepo:            userRepo,
		config:              cfg,
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

	// Ищем организацию по user ID
	company, err := as.organizationRepo.FindOrganizationByUserID(user.ID)
	if err != nil {
		return
	}

	// Проверяем, активен ли пользователь
	if !user.IsActive {
		err = errors.New("user is blocked")
		return
	}

	// do create token
	token, exp, err := jwt.
		NewJWT(as.config.Middleware.Jwt.Secret, as.config.Middleware.Jwt.Expiration).
		GenerateToken(jwt.JWTData{
			ID:        user.ID,
			Email:     user.Email,
			CompanyID: company.ID,
			Role:      string(user.Role),
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
		Role:     storage.RoleUser,
		IsActive: true,
	}

	user, err = as.userRepo.CreateUser(newUser)
	if err != nil {
		return
	}

	// Создание компании и связывание с пользователем
	if err = as.organizationService.CreateOrganization(user.ID, req.CompanyName, req.Phone); err != nil {
		return
	}

	res.ID = user.ID
	res.Email = user.Email
	res.Name = user.Name
	res.Company = user.Company
	res.Phone = user.Phone

	return
}

func (as *authService) ChangePassword(userID uint, req payload.ChangePasswordRequest) (res payload.ChangePasswordResponse, err error) {
	// Find user by ID
	user, err := as.userRepo.FindUserByID(int64(userID))
	if err != nil {
		return payload.ChangePasswordResponse{Success: false}, errors.New("user not found")
	}

	// Verify old password
	if !user.ComparePassword(req.OldPassword) {
		return payload.ChangePasswordResponse{Success: false}, errors.New("old password is incorrect")
	}

	// Hash new password
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return payload.ChangePasswordResponse{Success: false}, err
	}

	// Update password
	err = as.userRepo.UpdatePassword(userID, string(bcryptPassword))
	if err != nil {
		return payload.ChangePasswordResponse{Success: false}, err
	}

	return payload.ChangePasswordResponse{Success: true}, nil
}

// RequestPasswordReset handles a request to reset a password
func (as *authService) RequestPasswordReset(req payload.RequestPasswordResetRequest) (res payload.RequestPasswordResetResponse, err error) {
	// Find user by email
	user, err := as.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		// Don't reveal that the email doesn't exist for security reasons
		return payload.RequestPasswordResetResponse{
			Success: true,
			Message: "If your email is registered, you will receive a password reset code",
		}, nil
	}

	// Generate a random 6-digit code
	code := strconv.Itoa(100000 + rand.Intn(900000)) // 6-digit code

	// Set expiration time (e.g., 1 hour from now)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Create password reset code
	_, err = as.userRepo.CreatePasswordResetCode(user.ID, code, expiresAt)
	if err != nil {
		return payload.RequestPasswordResetResponse{Success: false}, err
	}

	// In a real application, you would send an email with the code
	// For this implementation, we'll just return success
	// TODO: Implement email sending

	return payload.RequestPasswordResetResponse{
		Success: true,
		Message: "If your email is registered, you will receive a password reset code",
	}, nil
}

// VerifyPasswordReset verifies a password reset code and sets a new password
func (as *authService) VerifyPasswordReset(req payload.VerifyPasswordResetRequest) (res payload.VerifyPasswordResetResponse, err error) {
	// Find user by email
	user, err := as.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return payload.VerifyPasswordResetResponse{Success: false}, errors.New("invalid email or code")
	}

	// Find password reset code
	passwordResetCode, err := as.userRepo.FindPasswordResetCodeByCode(req.Code)
	if err != nil {
		return payload.VerifyPasswordResetResponse{Success: false}, errors.New("invalid email or code")
	}

	// Check if the code belongs to the user
	if passwordResetCode.UserID != user.ID {
		return payload.VerifyPasswordResetResponse{Success: false}, errors.New("invalid email or code")
	}

	// Check if the code is valid
	if !passwordResetCode.IsValid() {
		return payload.VerifyPasswordResetResponse{Success: false}, errors.New("code is expired or already used")
	}

	// Hash new password
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return payload.VerifyPasswordResetResponse{Success: false}, err
	}

	// Update password
	err = as.userRepo.UpdatePassword(user.ID, string(bcryptPassword))
	if err != nil {
		return payload.VerifyPasswordResetResponse{Success: false}, err
	}

	// Mark code as used
	err = as.userRepo.MarkPasswordResetCodeAsUsed(passwordResetCode.ID)
	if err != nil {
		// Log this error but don't return it to the user
		// The password has been changed successfully
	}

	return payload.VerifyPasswordResetResponse{
		Success: true,
		Message: "Password has been reset successfully",
	}, nil
}
