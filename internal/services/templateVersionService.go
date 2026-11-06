package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

type TemplateVersionService struct {
	repo         *repositories.TemplateVersionRepository
	templateRepo *repositories.TemplateRepository
}

func NewTemplateVersionService(repo *repositories.TemplateVersionRepository, templateRepo *repositories.TemplateRepository) *TemplateVersionService {
	return &TemplateVersionService{repo: repo, templateRepo: templateRepo}
}

func (s *TemplateVersionService) Create(ctx context.Context, templateID string, req *dto.CreateTemplateVersionRequest) (*models.TemplateVersion, error) {
	templateID = strings.TrimSpace(templateID)
	version := strings.TrimSpace(req.Version)
	if templateID == "" || version == "" {
		return nil, errors.New("template_id dan version wajib diisi")
	}
	if len(version) > 20 {
		return nil, errors.New("version maksimal 20 karakter")
	}
	if _, err := s.templateRepo.GetByID(ctx, templateID); err != nil {
		return nil, errors.New("template tidak ditemukan")
	}
	if existing, _ := s.repo.GetByVersion(ctx, templateID, version); existing != nil {
		return nil, errors.New("version sudah digunakan pada template ini")
	}
	config := req.ConfigSchema
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	if !json.Valid(config) {
		return nil, errors.New("config_schema harus berupa JSON yang valid")
	}
	v := &models.TemplateVersion{TemplateID: templateID, Version: version, ConfigSchema: config}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	if req.IsActive != nil && *req.IsActive {
		if err := s.repo.SetActive(ctx, templateID, v.ID); err != nil {
			return nil, err
		}
		v.IsActive = true
	}
	return v, nil
}

func (s *TemplateVersionService) GetAll(ctx context.Context, templateID string) ([]models.TemplateVersion, error) {
	return s.repo.GetAll(ctx, templateID)
}

func (s *TemplateVersionService) GetByID(ctx context.Context, templateID, id string) (*models.TemplateVersion, error) {
	return s.repo.GetByID(ctx, templateID, id)
}

func (s *TemplateVersionService) Activate(ctx context.Context, templateID, id string) error {
	if _, err := s.repo.GetByID(ctx, templateID, id); err != nil {
		return err
	}
	return s.repo.SetActive(ctx, templateID, id)
}

func (s *TemplateVersionService) Delete(ctx context.Context, templateID, id string) error {
	v, err := s.repo.GetByID(ctx, templateID, id)
	if err != nil {
		return err
	}
	if v.IsActive {
		return errors.New("versi aktif tidak dapat dihapus; aktifkan versi lain terlebih dahulu")
	}
	return s.repo.Delete(ctx, templateID, id)
}
