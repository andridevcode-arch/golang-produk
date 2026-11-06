package services

import (
	"context"
	"errors"
	"strings"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*models.Category, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	existing, _ := s.repo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("slug kategori sudah digunakan")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cat := &models.Category{
		ParentID:    req.ParentID,
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Icon:        req.Icon,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    isActive,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context, onlyActive bool, asTree bool) ([]models.Category, error) {
	cats, err := s.repo.GetAll(ctx, onlyActive)
	if err != nil {
		return nil, err
	}

	if !asTree {
		return cats, nil
	}

	// Build category tree (hierarki parent-child)
	catMap := make(map[string]*models.Category)
	var rootCategories []models.Category

	for i := range cats {
		cats[i].Children = []models.Category{}
		catMap[cats[i].ID] = &cats[i]
	}

	for _, cat := range cats {
		if cat.ParentID != nil && *cat.ParentID != "" {
			if parent, exists := catMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, *catMap[cat.ID])
			} else {
				rootCategories = append(rootCategories, *catMap[cat.ID])
			}
		} else {
			rootCategories = append(rootCategories, *catMap[cat.ID])
		}
	}

	return rootCategories, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*models.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		cat.Name = strings.TrimSpace(req.Name)
	}

	if req.Slug != "" {
		newSlug := strings.TrimSpace(req.Slug)
		if newSlug != cat.Slug {
			existing, _ := s.repo.GetBySlug(ctx, newSlug)
			if existing != nil && existing.ID != cat.ID {
				return nil, errors.New("slug kategori sudah digunakan")
			}
			cat.Slug = newSlug
		}
	}

	if req.ParentID != nil {
		if *req.ParentID == cat.ID {
			return nil, errors.New("kategori tidak bisa menjadi parent untuk dirinya sendiri")
		}
		cat.ParentID = req.ParentID
	}

	if req.Icon != nil {
		cat.Icon = req.Icon
	}

	if req.Description != nil {
		cat.Description = req.Description
	}

	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}

	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
