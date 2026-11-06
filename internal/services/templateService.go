package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

type TemplateService struct {
	repo        *repositories.TemplateRepository
	productRepo *repositories.ProductRepository
}

func NewTemplateService(repo *repositories.TemplateRepository, productRepo *repositories.ProductRepository) *TemplateService {
	return &TemplateService{repo: repo, productRepo: productRepo}
}

func (s *TemplateService) CreateTemplate(ctx context.Context, req *dto.CreateTemplateRequest) (*models.Template, error) {
	productID := strings.TrimSpace(req.ProductID)
	name := strings.TrimSpace(req.Name)
	previewURL := strings.TrimSpace(req.PreviewURL)
	if productID == "" || name == "" || previewURL == "" {
		return nil, errors.New("product_id, name, dan preview_url wajib diisi")
	}
	if _, err := s.productRepo.GetByID(ctx, productID); err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	config := req.ConfigSchema
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	if !json.Valid(config) {
		return nil, errors.New("config_schema harus berupa JSON yang valid")
	}
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}
	if existing, _ := s.repo.GetBySlug(ctx, productID, slug); existing != nil {
		return nil, errors.New("slug template sudah digunakan pada produk ini")
	}
	premium, active := false, true
	if req.IsPremium != nil {
		premium = *req.IsPremium
	}
	if req.IsActive != nil {
		active = *req.IsActive
	}
	t := &models.Template{ProductID: productID, Name: name, Slug: slug, PreviewURL: previewURL, ThumbnailURL: req.ThumbnailURL, IsPremium: premium, IsActive: active, SortOrder: req.SortOrder, ConfigSchema: config}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) GetAllTemplates(ctx context.Context, productID string, premium *bool, onlyActive bool) ([]models.Template, error) {
	return s.repo.GetAll(ctx, productID, premium, onlyActive)
}

func (s *TemplateService) GetTemplateByID(ctx context.Context, id string) (*models.Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TemplateService) UpdateTemplate(ctx context.Context, id string, req *dto.UpdateTemplateRequest) (*models.Template, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		t.Name = strings.TrimSpace(*req.Name)
	}
	if req.Slug != nil {
		slug := strings.TrimSpace(*req.Slug)
		if slug == "" {
			return nil, errors.New("slug tidak boleh kosong")
		}
		if existing, _ := s.repo.GetBySlug(ctx, t.ProductID, slug); existing != nil && existing.ID != t.ID {
			return nil, errors.New("slug template sudah digunakan pada produk ini")
		}
		t.Slug = slug
	}
	if req.PreviewURL != nil {
		if strings.TrimSpace(*req.PreviewURL) == "" {
			return nil, errors.New("preview_url tidak boleh kosong")
		}
		t.PreviewURL = strings.TrimSpace(*req.PreviewURL)
	}
	if req.ThumbnailURL != nil {
		t.ThumbnailURL = req.ThumbnailURL
	}
	if req.IsPremium != nil {
		t.IsPremium = *req.IsPremium
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		t.SortOrder = *req.SortOrder
	}
	if req.ConfigSchema != nil {
		if !json.Valid(req.ConfigSchema) {
			return nil, errors.New("config_schema harus berupa JSON yang valid")
		}
		t.ConfigSchema = req.ConfigSchema
	}
	if t.ConfigSchema == nil {
		t.ConfigSchema = json.RawMessage(`{}`)
	}
	t.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) DeleteTemplate(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
