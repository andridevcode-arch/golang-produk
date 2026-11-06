package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-produk/internal/dto"
	"go-produk/internal/middlewares"
	"go-produk/internal/repositories"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

type PublishedPageHandler struct {
	service *services.PublishedPageService
}

func NewPublishedPageHandler(service *services.PublishedPageService) *PublishedPageHandler {
	return &PublishedPageHandler{service: service}
}
func (h *PublishedPageHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var req dto.UpsertPublishedPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	page, err := h.service.Upsert(r.Context(), r.PathValue("id"), middlewares.GetUserID(r), &req)
	h.respondOwner(w, page, err, http.StatusCreated)
}
func (h *PublishedPageHandler) GetOwner(w http.ResponseWriter, r *http.Request) {
	page, err := h.service.GetOwner(r.Context(), r.PathValue("id"), middlewares.GetUserID(r))
	h.respondOwner(w, page, err, http.StatusOK)
}
func (h *PublishedPageHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePublishedPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	page, err := h.service.Update(r.Context(), r.PathValue("id"), middlewares.GetUserID(r), &req)
	h.respondOwner(w, page, err, http.StatusOK)
}
func (h *PublishedPageHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.changeState(w, r, true)
}
func (h *PublishedPageHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
	h.changeState(w, r, false)
}
func (h *PublishedPageHandler) changeState(w http.ResponseWriter, r *http.Request, active bool) {
	var err error
	if active {
		err = h.service.Publish(r.Context(), r.PathValue("id"), middlewares.GetUserID(r))
	} else {
		err = h.service.Unpublish(r.Context(), r.PathValue("id"), middlewares.GetUserID(r))
	}
	if errors.Is(err, repositories.ErrPublishedPageNotFound) || errors.Is(err, repositories.ErrUserProjectNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, nil)
}
func (h *PublishedPageHandler) Public(w http.ResponseWriter, r *http.Request) {
	page, err := h.service.GetPublic(r.Context(), r.PathValue("username"), r.PathValue("slug"), "", false)
	h.respondPublic(w, page, err)
}
func (h *PublishedPageHandler) Preview(w http.ResponseWriter, r *http.Request) {
	page, err := h.service.GetPublic(r.Context(), r.PathValue("username"), r.PathValue("slug"), r.URL.Query().Get("preview_token"), true)
	h.respondPublic(w, page, err)
}
func (h *PublishedPageHandler) respondOwner(w http.ResponseWriter, page any, err error, status int) {
	if errors.Is(err, repositories.ErrPublishedPageNotFound) || errors.Is(err, repositories.ErrUserProjectNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	if status == http.StatusCreated {
		response.RespondCreated(w, page)
		return
	}
	response.RespondSuccess(w, page)
}
func (h *PublishedPageHandler) respondPublic(w http.ResponseWriter, page any, err error) {
	if errors.Is(err, repositories.ErrPublishedPageNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, page)
}
