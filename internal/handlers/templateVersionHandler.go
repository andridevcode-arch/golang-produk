package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-produk/internal/dto"
	"go-produk/internal/repositories"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

type TemplateVersionHandler struct {
	service *services.TemplateVersionService
}

func NewTemplateVersionHandler(service *services.TemplateVersionService) *TemplateVersionHandler {
	return &TemplateVersionHandler{service: service}
}

func (h *TemplateVersionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll(r.Context(), r.PathValue("template_id"))
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccess(w, items)
}

func (h *TemplateVersionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.GetByID(r.Context(), r.PathValue("template_id"), r.PathValue("version_id"))
	if errors.Is(err, repositories.ErrTemplateVersionNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccess(w, item)
}

func (h *TemplateVersionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTemplateVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	item, err := h.service.Create(r.Context(), r.PathValue("template_id"), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondCreated(w, item)
}

func (h *TemplateVersionHandler) Activate(w http.ResponseWriter, r *http.Request) {
	err := h.service.Activate(r.Context(), r.PathValue("template_id"), r.PathValue("version_id"))
	if errors.Is(err, repositories.ErrTemplateVersionNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccessWithMessage(w, "versi template berhasil diaktifkan", nil)
}

func (h *TemplateVersionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(r.Context(), r.PathValue("template_id"), r.PathValue("version_id"))
	if errors.Is(err, repositories.ErrTemplateVersionNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccessWithMessage(w, "versi template berhasil dihapus", nil)
}
