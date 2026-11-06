package dto

type CreateUserProjectRequest struct {
	TemplateID        string  `json:"template_id"`
	TemplateVersionID *string `json:"template_version_id"`
	Title             string  `json:"title"`
	Description       *string `json:"description"`
}

type UpdateUserProjectRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UpdateUserProjectStatusRequest struct {
	Status string `json:"status"`
}
