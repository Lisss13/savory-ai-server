package service

import (
	"savory-ai-server/app/module/organization/payload"
	"savory-ai-server/app/module/organization/repository"
	"savory-ai-server/app/storage"
)

type languageService struct {
	langRepo repository.LanguageRepository
}

type LanguageService interface {
	GetAllLanguages() (*payload.LanguagesResp, error)
	GetLanguageByID(id uint) (*payload.LanguageResp, error)
	CreateLanguage(req *payload.CreateLanguageReq) (*payload.LanguageResp, error)
	UpdateLanguage(id uint, req *payload.UpdateLanguageReq) (*payload.LanguageResp, error)
	DeleteLanguage(id uint) error
	
	// Organization language operations
	GetLanguagesByOrganizationID(orgID uint) (*payload.OrganizationLanguagesResp, error)
	AddLanguageToOrganization(orgID uint, req *payload.AddLanguageToOrgReq) error
	RemoveLanguageFromOrganization(orgID uint, req *payload.RemoveLanguageFromOrgReq) error
}

func NewLanguageService(langRepo repository.LanguageRepository) LanguageService {
	return &languageService{
		langRepo: langRepo,
	}
}

func (ls *languageService) GetAllLanguages() (*payload.LanguagesResp, error) {
	languages, err := ls.langRepo.FindAllLanguages()
	if err != nil {
		return nil, err
	}

	resp := &payload.LanguagesResp{
		Languages: make([]payload.LanguageResp, 0, len(languages)),
	}

	for _, lang := range languages {
		resp.Languages = append(resp.Languages, *mapLanguageToResponse(lang))
	}

	return resp, nil
}

func (ls *languageService) GetLanguageByID(id uint) (*payload.LanguageResp, error) {
	language, err := ls.langRepo.FindLanguageByID(id)
	if err != nil {
		return nil, err
	}

	return mapLanguageToResponse(language), nil
}

func (ls *languageService) CreateLanguage(req *payload.CreateLanguageReq) (*payload.LanguageResp, error) {
	// Check if language with the same code already exists
	existingLang, err := ls.langRepo.FindLanguageByCode(req.Code)
	if err == nil && existingLang != nil {
		// Language with this code already exists
		return mapLanguageToResponse(existingLang), nil
	}

	// Create new language
	language := &storage.Language{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	createdLang, err := ls.langRepo.CreateLanguage(language)
	if err != nil {
		return nil, err
	}

	return mapLanguageToResponse(createdLang), nil
}

func (ls *languageService) UpdateLanguage(id uint, req *payload.UpdateLanguageReq) (*payload.LanguageResp, error) {
	// Check if language exists
	existingLang, err := ls.langRepo.FindLanguageByID(id)
	if err != nil {
		return nil, err
	}

	// Update language fields if provided
	if req.Code != "" {
		existingLang.Code = req.Code
	}
	if req.Name != "" {
		existingLang.Name = req.Name
	}
	if req.Description != "" {
		existingLang.Description = req.Description
	}

	updatedLang, err := ls.langRepo.UpdateLanguage(existingLang)
	if err != nil {
		return nil, err
	}

	return mapLanguageToResponse(updatedLang), nil
}

func (ls *languageService) DeleteLanguage(id uint) error {
	return ls.langRepo.DeleteLanguage(id)
}

func (ls *languageService) GetLanguagesByOrganizationID(orgID uint) (*payload.OrganizationLanguagesResp, error) {
	languages, err := ls.langRepo.FindLanguagesByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}

	resp := &payload.OrganizationLanguagesResp{
		OrganizationID: orgID,
		Languages:      make([]payload.LanguageResp, 0, len(languages)),
	}

	for _, lang := range languages {
		resp.Languages = append(resp.Languages, *mapLanguageToResponse(lang))
	}

	return resp, nil
}

func (ls *languageService) AddLanguageToOrganization(orgID uint, req *payload.AddLanguageToOrgReq) error {
	return ls.langRepo.AddLanguageToOrganization(orgID, req.LanguageID)
}

func (ls *languageService) RemoveLanguageFromOrganization(orgID uint, req *payload.RemoveLanguageFromOrgReq) error {
	return ls.langRepo.RemoveLanguageFromOrganization(orgID, req.LanguageID)
}

// Helper function to map storage.Language to payload.LanguageResp
func mapLanguageToResponse(lang *storage.Language) *payload.LanguageResp {
	return &payload.LanguageResp{
		ID:          lang.ID,
		CreatedAt:   lang.CreatedAt,
		Code:        lang.Code,
		Name:        lang.Name,
		Description: lang.Description,
	}
}