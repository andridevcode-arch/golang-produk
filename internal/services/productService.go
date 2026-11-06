package services

import (
	"context"
	"errors"
	"strings"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*models.Product, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("nama produk wajib diisi")
	}

	validAccessTypes := map[string]bool{"free": true, "premium": true, "one_time": true}
	accessType := strings.TrimSpace(req.AccessType)
	if accessType == "" {
		accessType = "free"
	}
	if !validAccessTypes[accessType] {
		return nil, errors.New("access_type harus salah satu dari: free, premium, one_time")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	existing, _ := s.repo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("slug produk sudah digunakan")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	p := &models.Product{
		CategoryID:  strings.TrimSpace(req.CategoryID),
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		Price:       req.Price,
		AccessType:  accessType,
		IsActive:    isActive,
		SortOrder:   req.SortOrder,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context, categoryID, accessType string, onlyActive bool) ([]models.Product, error) {
	return s.repo.GetAll(ctx, categoryID, accessType, onlyActive)
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		p.Name = strings.TrimSpace(*req.Name)
	}

	if req.Slug != nil {
		newSlug := strings.TrimSpace(*req.Slug)
		if newSlug != p.Slug {
			existing, _ := s.repo.GetBySlug(ctx, newSlug)
			if existing != nil && existing.ID != p.ID {
				return nil, errors.New("slug produk sudah digunakan")
			}
			p.Slug = newSlug
		}
	}

	if req.Description != nil {
		p.Description = req.Description
	}

	if req.Price != nil {
		p.Price = req.Price
	}

	if req.AccessType != nil {
		validAccessTypes := map[string]bool{"free": true, "premium": true, "one_time": true}
		if !validAccessTypes[*req.AccessType] {
			return nil, errors.New("access_type harus salah satu dari: free, premium, one_time")
		}
		p.AccessType = *req.AccessType
	}

	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	if req.SortOrder != nil {
		p.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
