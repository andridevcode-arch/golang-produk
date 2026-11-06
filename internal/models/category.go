package models

import "time"

type Category struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parent_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Icon        *string    `json:"icon"`
	Description *string    `json:"description"`
	SortOrder   int        `json:"sort_order"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Children    []Category `json:"children,omitempty"`
}
