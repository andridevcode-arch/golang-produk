package response

import "net/http"

// PaginationMeta metadata pagination
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// PaginationResponse struktur response pagination
type PaginationResponse struct {
	Success    bool           `json:"success"`
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Meta       PaginationMeta `json:"meta"`
}

// NewPaginationResponse membuat response pagination baru
func NewPaginationResponse(data interface{}, currentPage, perPage, totalItems int) *PaginationResponse {
	totalPages := totalItems / perPage
	if totalItems%perPage > 0 {
		totalPages++
	}

	return &PaginationResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "sukses",
		Data:       data,
		Meta: PaginationMeta{
			CurrentPage: currentPage,
			PerPage:     perPage,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}
}

// RespondPaginated merespon dengan data yang sudah dipaginasi
func RespondPaginated(w http.ResponseWriter, data interface{}, currentPage, perPage, totalItems int) {
	RespondJSON(w, http.StatusOK, NewPaginationResponse(data, currentPage, perPage, totalItems))
}
