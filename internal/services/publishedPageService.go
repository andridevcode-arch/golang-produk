package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type PublishedPageService struct {
	repo *repositories.PublishedPageRepository
}

func NewPublishedPageService(repo *repositories.PublishedPageRepository) *PublishedPageService {
	return &PublishedPageService{repo: repo}
}

func (s *PublishedPageService) Upsert(ctx context.Context, projectID, userID string, req *dto.UpsertPublishedPageRequest) (*models.PublishedPage, error) {
	slug, err := normalizeSlug(req.Slug)
	if err != nil {
		return nil, err
	}
	project, err := s.repo.GetOwnerProject(ctx, strings.TrimSpace(projectID), strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	if !project.IsActive || project.Status == "archived" {
		return nil, errors.New("proyek tidak dapat memiliki halaman publik")
	}
	page, err := s.repo.GetByProjectAndUserID(ctx, project.ID, project.UserID)
	if err != nil && !errors.Is(err, repositories.ErrPublishedPageNotFound) {
		return nil, err
	}
	username, err := s.ownerUsername(ctx, project.UserID)
	if err != nil {
		return nil, err
	}
	allowPreview := true
	if page != nil {
		allowPreview = page.AllowPreview
	}
	if req.AllowPreview != nil {
		allowPreview = *req.AllowPreview
	}
	if page == nil {
		page = &models.PublishedPage{ProjectID: project.ID, UserID: project.UserID, PreviewToken: randomToken()}
	}
	page.Username, page.Slug, page.CustomPath, page.AllowPreview = username, slug, "/"+username+"/"+slug, allowPreview
	if err := s.repo.Upsert(ctx, page); err != nil {
		return nil, err
	}
	page.Project = project
	return page, nil
}

func (s *PublishedPageService) GetOwner(ctx context.Context, projectID, userID string) (*models.PublishedPage, error) {
	page, err := s.repo.GetByProjectAndUserID(ctx, strings.TrimSpace(projectID), strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	page.Project, err = s.repo.GetOwnerProject(ctx, page.ProjectID, page.UserID)
	return page, err
}

func (s *PublishedPageService) Update(ctx context.Context, projectID, userID string, req *dto.UpdatePublishedPageRequest) (*models.PublishedPage, error) {
	page, err := s.repo.GetByProjectAndUserID(ctx, strings.TrimSpace(projectID), strings.TrimSpace(userID))
	if err != nil {
		return nil, err
	}
	if req.Slug != nil {
		page.Slug, err = normalizeSlug(*req.Slug)
		if err != nil {
			return nil, err
		}
	}
	if req.AllowPreview != nil {
		page.AllowPreview = *req.AllowPreview
	}
	username, err := s.ownerUsername(ctx, page.UserID)
	if err != nil {
		return nil, err
	}
	page.Username, page.CustomPath = username, "/"+username+"/"+page.Slug
	if err := s.repo.Upsert(ctx, page); err != nil {
		return nil, err
	}
	page.Project, err = s.repo.GetOwnerProject(ctx, page.ProjectID, page.UserID)
	return page, err
}

func (s *PublishedPageService) Publish(ctx context.Context, projectID, userID string) error {
	return s.repo.Publish(ctx, strings.TrimSpace(projectID), strings.TrimSpace(userID), true)
}
func (s *PublishedPageService) Unpublish(ctx context.Context, projectID, userID string) error {
	return s.repo.Publish(ctx, strings.TrimSpace(projectID), strings.TrimSpace(userID), false)
}

func (s *PublishedPageService) GetPublic(ctx context.Context, username, slug, previewToken string, preview bool) (*dto.PublicPublishedPageResponse, error) {
	page, err := s.repo.GetPublic(ctx, strings.TrimSpace(username), strings.TrimSpace(slug), preview)
	if err != nil {
		return nil, err
	}
	if preview && !page.AllowPreview && previewToken != page.PreviewToken {
		return nil, errors.New("preview token tidak valid")
	}
	project, err := s.repo.GetOwnerProject(ctx, page.ProjectID, page.UserID)
	if err != nil {
		return nil, err
	}
	if !preview {
		if err := s.repo.IncrementViews(ctx, page.ID); err != nil {
			return nil, err
		}
		page.ViewCount++
	}
	return &dto.PublicPublishedPageResponse{ID: page.ID, ProjectID: page.ProjectID, Username: page.Username, Slug: page.Slug, CustomPath: page.CustomPath, AllowPreview: page.AllowPreview, IsActive: page.IsActive, ViewCount: page.ViewCount, PublishedAt: page.PublishedAt, CreatedAt: page.CreatedAt, UpdatedAt: page.UpdatedAt, Project: project}, nil
}

func (s *PublishedPageService) ownerUsername(ctx context.Context, userID string) (string, error) {
	username, err := s.repo.GetUsername(ctx, userID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(username) == "" {
		return "", errors.New("username pemilik belum tersedia")
	}
	return username, nil
}

func normalizeSlug(value string) (string, error) {
	slug := strings.ToLower(strings.TrimSpace(value))
	if len(slug) == 0 || len(slug) > 150 || !slugPattern.MatchString(slug) {
		return "", errors.New("slug harus maksimal 150 karakter, hanya huruf kecil, angka, dan tanda hubung")
	}
	return slug, nil
}
func randomToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
