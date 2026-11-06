package models

import (
	"encoding/json"
	"time"
)

type Template struct {
	ID           string          `json:"id"`
	ProductID    string          `json:"product_id"`
	Products     []Product       `json:"products,omitempty"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	PreviewURL   string          `json:"preview_url"`
	ThumbnailURL *string         `json:"thumbnail_url,omitempty"`
	IsPremium    bool            `json:"is_premium"`
	IsActive     bool            `json:"is_active"`
	SortOrder    int             `json:"sort_order"`
	ConfigSchema json.RawMessage `json:"config_schema"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
