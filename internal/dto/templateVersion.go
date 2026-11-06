package dto

import "encoding/json"

type CreateTemplateVersionRequest struct {
	Version      string          `json:"version"`
	ConfigSchema json.RawMessage `json:"config_schema"`
	IsActive     *bool           `json:"is_active"`
}
