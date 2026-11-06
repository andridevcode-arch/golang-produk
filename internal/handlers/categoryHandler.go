package handlers

import (
	"encoding/json"
	"net/http"

	"go-produk/internal/dto"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll menangani GET /api/v1/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	onlyActive := r.URL.Query().Get("active") == "true"
	asTree := r.URL.Query().Get("tree") == "true"

	categories, err := h.service.GetAllCategories(r.Context(), onlyActive, asTree)
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}

	response.RespondSuccess(w, categories)
}

// GetByID menangani GET /api/v1/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id kategori diperlukan")
		return
	}

	cat, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		response.RespondNotFound(w, err.Error())
		return
	}

	response.RespondSuccess(w, cat)
}

// Create menangani POST /api/v1/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}

	cat, err := h.service.CreateCategory(r.Context(), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondCreated(w, cat)
}

// Update menangani PUT /api/v1/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id kategori diperlukan")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}

	cat, err := h.service.UpdateCategory(r.Context(), id, &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccess(w, cat)
}

// Delete menangani DELETE /api/v1/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id kategori diperlukan")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccess(w, map[string]string{"message": "kategori berhasil dihapus"})
}
