package models

import "time"

type UserProject struct {
	ID                string           `json:"id"`
	UserID            string           `json:"user_id"`
	TemplateID        string           `json:"template_id"`
	TemplateVersionID string           `json:"template_version_id"`
	Template          *Template        `json:"template,omitempty"`
	TemplateVersion   *TemplateVersion `json:"template_version,omitempty"`
	Title             string           `json:"title"`
	Description       *string          `json:"description,omitempty"`
	Status            string           `json:"status"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
