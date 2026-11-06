package models

import "time"

type Product struct {
	ID          string     `json:"id"`
	CategoryID  string     `json:"category_id"`
	Categories  []Category `json:"categories,omitempty"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	AccessType  string     `json:"access_type"`
	IsActive    bool       `json:"is_active"`
	IsDeleted   bool       `json:"is_deleted"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
