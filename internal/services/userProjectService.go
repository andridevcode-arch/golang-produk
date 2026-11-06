package services

import (
	"context"
	"errors"
	"strings"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

type UserProjectService struct {
	repo         *repositories.UserProjectRepository
	templateRepo *repositories.TemplateRepository
}

func NewUserProjectService(repo *repositories.UserProjectRepository, templateRepo *repositories.TemplateRepository) *UserProjectService {
	return &UserProjectService{repo: repo, templateRepo: templateRepo}
}

func (s *UserProjectService) Create(ctx context.Context, userID string, req *dto.CreateUserProjectRequest) (*models.UserProject, error) {
	userID = strings.TrimSpace(userID)
	templateID := strings.TrimSpace(req.TemplateID)
	title := strings.TrimSpace(req.Title)
	if userID == "" || templateID == "" || title == "" {
		return nil, errors.New("template_id dan title wajib diisi")
	}
	if len(title) > 150 {
		return nil, errors.New("title maksimal 150 karakter")
	}

	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, errors.New("template tidak ditemukan")
	}
	if !template.IsActive {
		return nil, errors.New("template tidak aktif")
	}

	versionID := ""
	if req.TemplateVersionID != nil {
		versionID = strings.TrimSpace(*req.TemplateVersionID)
	}
	if versionID == "" {
		versionID, err = s.repo.GetActiveTemplateVersionID(ctx, templateID)
		if errors.Is(err, repositories.ErrTemplateVersionNotFound) {
			return nil, errors.New("template belum memiliki versi aktif")
		}
		if err != nil {
			return nil, err
		}
	} else {
		owned, err := s.repo.IsTemplateVersionOwnedByTemplate(ctx, templateID, versionID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, errors.New("versi template tidak sesuai dengan template yang dipilih")
		}
	}

	project := &models.UserProject{
		UserID:            userID,
		TemplateID:        templateID,
		TemplateVersionID: versionID,
		Title:             title,
		Description:       req.Description,
	}
	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}
	return s.repo.GetByIDAndUserID(ctx, project.ID, userID)
}

func (s *UserProjectService) List(ctx context.Context, userID, status string) ([]models.UserProject, error) {
	status = strings.TrimSpace(strings.ToLower(status))
	if status != "" && !isValidProjectStatus(status) {
		return nil, errors.New("status harus draft, published, atau archived")
	}
	return s.repo.ListByUserID(ctx, userID, status)
}

func (s *UserProjectService) GetByID(ctx context.Context, id, userID string) (*models.UserProject, error) {
	return s.repo.GetByIDAndUserID(ctx, strings.TrimSpace(id), strings.TrimSpace(userID))
}

func (s *UserProjectService) Update(ctx context.Context, id, userID string, req *dto.UpdateUserProjectRequest) (*models.UserProject, error) {
	project, err := s.repo.GetByIDAndUserID(ctx, strings.TrimSpace(id), strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, errors.New("title tidak boleh kosong")
		}
		if len(title) > 150 {
			return nil, errors.New("title maksimal 150 karakter")
		}
		project.Title = title
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if err := s.repo.UpdateByIDAndUserID(ctx, project); err != nil {
		return nil, err
	}
	return s.repo.GetByIDAndUserID(ctx, project.ID, project.UserID)
}

func (s *UserProjectService) UpdateStatus(ctx context.Context, id, userID string, req *dto.UpdateUserProjectStatusRequest) error {
	status := strings.TrimSpace(strings.ToLower(req.Status))
	if status == "published" {
		return errors.New("status published hanya dapat diubah melalui endpoint halaman publik")
	}
	if status != "draft" && status != "archived" {
		return errors.New("status harus draft atau archived")
	}
	if status == "archived" {
		return s.repo.ArchiveByIDAndUserID(ctx, strings.TrimSpace(id), strings.TrimSpace(userID))
	}
	return s.repo.SetStatusByIDAndUserID(ctx, strings.TrimSpace(id), strings.TrimSpace(userID), status)
}

func (s *UserProjectService) Delete(ctx context.Context, id, userID string) error {
	return s.repo.ArchiveByIDAndUserID(ctx, strings.TrimSpace(id), strings.TrimSpace(userID))
}

func isValidProjectStatus(status string) bool {
	return status == "draft" || status == "published" || status == "archived"
}
