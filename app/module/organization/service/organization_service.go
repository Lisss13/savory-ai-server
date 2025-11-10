package service

import (
	"savory-ai-server/app/module/organization/payload"
	"savory-ai-server/app/module/organization/repository"
	"savory-ai-server/app/storage"
)

type organizationService struct {
	orgRepo repository.OrganizationRepository
}

type OrganizationService interface {
	GetOrganizationByID(id uint) (*payload.OrganizationResp, error)
	GetAllOrganizations() (*payload.OrganizationsResp, error)
	CreateOrganization(adminID uint, name, phone string) error
	UpdateOrganization(id uint, req *payload.UpdateOrganizationReq) (*payload.OrganizationResp, error)
	AddUserToOrganization(orgID, adminID uint) error
	RemoveUserFromOrganization(orgID uint, req *payload.RemoveUserFromOrgReq) error
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}

func (os *organizationService) GetOrganizationByID(id uint) (*payload.OrganizationResp, error) {
	org, err := os.orgRepo.FindOrganizationByID(id)
	if err != nil {
		return nil, err
	}

	return mapOrganizationToResponse(org), nil
}

func (os *organizationService) GetAllOrganizations() (*payload.OrganizationsResp, error) {
	orgs, err := os.orgRepo.FindAllOrganizations()
	if err != nil {
		return nil, err
	}

	resp := &payload.OrganizationsResp{
		Organizations: make([]payload.OrganizationResp, 0, len(orgs)),
	}

	for _, org := range orgs {
		resp.Organizations = append(resp.Organizations, *mapOrganizationToResponse(org))
	}

	return resp, nil
}

// CreateOrganization creates a new organization and returns the complete organization with relationships
func (os *organizationService) CreateOrganization(adminID uint, name, phone string) error {
	org := &storage.Organization{
		Name:    name,
		Phone:   phone,
		AdminID: adminID,
	}

	orgStore, err := os.orgRepo.CreateOrganization(org)
	if err != nil {
		return err
	}

	if err = os.AddUserToOrganization(orgStore.ID, adminID); err != nil {
		return err
	}

	return nil
}

func (os *organizationService) UpdateOrganization(id uint, req *payload.UpdateOrganizationReq) (*payload.OrganizationResp, error) {
	// Check if organization exists
	existingOrg, err := os.orgRepo.FindOrganizationByID(id)
	if err != nil {
		return nil, err
	}

	// Update organization fields
	existingOrg.Name = req.Name
	existingOrg.Phone = req.Phone

	updatedOrg, err := os.orgRepo.UpdateOrganization(existingOrg)
	if err != nil {
		return nil, err
	}

	return mapOrganizationToResponse(updatedOrg), nil
}

func (os *organizationService) AddUserToOrganization(orgID, userID uint) error {
	return os.orgRepo.AddUserToOrganization(orgID, userID)
}

func (os *organizationService) RemoveUserFromOrganization(orgID uint, req *payload.RemoveUserFromOrgReq) error {
	return os.orgRepo.RemoveUserFromOrganization(orgID, req.UserID)
}

// Helper function to map storage.Organization to payload.OrganizationResp
func mapOrganizationToResponse(org *storage.Organization) *payload.OrganizationResp {
	resp := &payload.OrganizationResp{
		ID:        org.ID,
		CreatedAt: org.CreatedAt,
		Name:      org.Name,
		Phone:     org.Phone,
		Admin: payload.UserInOrgResp{
			ID:    org.Admin.ID,
			Name:  org.Admin.Name,
			Email: org.Admin.Email,
			Phone: org.Admin.Phone,
		},
	}

	if len(org.Users) > 0 {
		resp.Users = make([]payload.UserInOrgResp, 0, len(org.Users))
		for _, user := range org.Users {
			resp.Users = append(resp.Users, payload.UserInOrgResp{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Phone: user.Phone,
			})
		}
	}

	return resp
}
