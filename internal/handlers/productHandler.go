package handlers

import (
	"encoding/json"
	"net/http"

	"go-produk/internal/dto"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetAll menangani GET /api/v1/products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	accessType := r.URL.Query().Get("access_type")
	onlyActive := r.URL.Query().Get("active") == "true"

	products, err := h.service.GetAllProducts(r.Context(), categoryID, accessType, onlyActive)
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}

	response.RespondSuccess(w, products)
}

// GetByID menangani GET /api/v1/products/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id produk diperlukan")
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		response.RespondNotFound(w, err.Error())
		return
	}

	response.RespondSuccess(w, product)
}

// Create menangani POST /api/v1/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}

	product, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondCreated(w, product)
}

// Update menangani PUT /api/v1/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id produk diperlukan")
		return
	}

	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccess(w, product)
}

// Delete menangani DELETE /api/v1/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.RespondBadRequest(w, "id produk diperlukan")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccess(w, map[string]string{"message": "produk berhasil dihapus"})
}
