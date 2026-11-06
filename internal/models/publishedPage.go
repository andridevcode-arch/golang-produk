package models

import "time"

type PublishedPage struct {
	ID           string       `json:"id"`
	ProjectID    string       `json:"project_id"`
	UserID       string       `json:"user_id"`
	Username     string       `json:"username"`
	Slug         string       `json:"slug"`
	CustomPath   string       `json:"custom_path"`
	PreviewToken string       `json:"preview_token,omitempty"`
	AllowPreview bool         `json:"allow_preview"`
	IsActive     bool         `json:"is_active"`
	ViewCount    int64        `json:"view_count"`
	PublishedAt  *time.Time   `json:"published_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Project      *UserProject `json:"project,omitempty"`
}
