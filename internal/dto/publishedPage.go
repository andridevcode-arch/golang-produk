package dto

import (
	"go-produk/internal/models"
	"time"
)

type UpsertPublishedPageRequest struct {
	Slug         string `json:"slug"`
	AllowPreview *bool  `json:"allow_preview"`
}
type UpdatePublishedPageRequest struct {
	Slug         *string `json:"slug"`
	AllowPreview *bool   `json:"allow_preview"`
}
type PublicPublishedPageResponse struct {
	ID           string              `json:"id"`
	ProjectID    string              `json:"project_id"`
	Username     string              `json:"username"`
	Slug         string              `json:"slug"`
	CustomPath   string              `json:"custom_path"`
	AllowPreview bool                `json:"allow_preview"`
	IsActive     bool                `json:"is_active"`
	ViewCount    int64               `json:"view_count"`
	PublishedAt  *time.Time          `json:"published_at,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	Project      *models.UserProject `json:"project"`
}
