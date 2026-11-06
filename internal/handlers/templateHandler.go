package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-produk/internal/dto"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

type TemplateHandler struct{ service *services.TemplateService }

func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (h *TemplateHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var premium *bool
	if v := q.Get("is_premium"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			response.RespondBadRequest(w, "is_premium harus boolean")
			return
		}
		premium = &b
	}
	items, err := h.service.GetAllTemplates(r.Context(), q.Get("product_id"), premium, q.Get("active") == "true")
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccess(w, items)
}

func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.GetTemplateByID(r.Context(), r.PathValue("id"))
	if err != nil {
		response.RespondNotFound(w, err.Error())
		return
	}
	response.RespondSuccess(w, item)
}

func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	item, err := h.service.CreateTemplate(r.Context(), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondCreated(w, item)
}

func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	item, err := h.service.UpdateTemplate(r.Context(), r.PathValue("id"), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, item)
}

func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteTemplate(r.Context(), r.PathValue("id")); err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, map[string]string{"message": "template berhasil dihapus"})
}
