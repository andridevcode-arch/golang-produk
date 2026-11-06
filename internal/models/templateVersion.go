package models

import (
	"encoding/json"
	"time"
)

type TemplateVersion struct {
	ID           string          `json:"id"`
	TemplateID   string          `json:"template_id"`
	Template     *Template       `json:"template,omitempty"`
	Version      string          `json:"version"`
	ConfigSchema json.RawMessage `json:"config_schema"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
}
