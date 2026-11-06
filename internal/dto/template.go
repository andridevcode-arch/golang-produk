package dto

import "encoding/json"

type CreateTemplateRequest struct {
	ProductID    string          `json:"product_id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	PreviewURL   string          `json:"preview_url"`
	ThumbnailURL *string         `json:"thumbnail_url"`
	IsPremium    *bool           `json:"is_premium"`
	IsActive     *bool           `json:"is_active"`
	SortOrder    int             `json:"sort_order"`
	ConfigSchema json.RawMessage `json:"config_schema"`
}

type UpdateTemplateRequest struct {
	Name         *string         `json:"name"`
	Slug         *string         `json:"slug"`
	PreviewURL   *string         `json:"preview_url"`
	ThumbnailURL *string         `json:"thumbnail_url"`
	IsPremium    *bool           `json:"is_premium"`
	IsActive     *bool           `json:"is_active"`
	SortOrder    *int            `json:"sort_order"`
	ConfigSchema json.RawMessage `json:"config_schema"`
}
