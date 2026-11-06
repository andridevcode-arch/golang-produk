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

type UserProjectHandler struct {
	service *services.UserProjectService
}

func NewUserProjectHandler(service *services.UserProjectService) *UserProjectHandler {
	return &UserProjectHandler{service: service}
}

func (h *UserProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	project, err := h.service.Create(r.Context(), middlewares.GetUserID(r), &req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondCreated(w, project)
}

func (h *UserProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.List(r.Context(), middlewares.GetUserID(r), r.URL.Query().Get("status"))
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, projects)
}

func (h *UserProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	project, err := h.service.GetByID(r.Context(), r.PathValue("id"), middlewares.GetUserID(r))
	if errors.Is(err, repositories.ErrUserProjectNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccess(w, project)
}

func (h *UserProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	project, err := h.service.Update(r.Context(), r.PathValue("id"), middlewares.GetUserID(r), &req)
	if errors.Is(err, repositories.ErrUserProjectNotFound) {
		response.RespondNotFound(w, err.Error())
		return
	}
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccess(w, project)
}

func (h *UserProjectHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserProjectStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request body tidak valid")
		return
	}
	if err := h.service.UpdateStatus(r.Context(), r.PathValue("id"), middlewares.GetUserID(r), &req); err != nil {
		if errors.Is(err, repositories.ErrUserProjectNotFound) {
			response.RespondNotFound(w, err.Error())
			return
		}
		response.RespondBadRequest(w, err.Error())
		return
	}
	response.RespondSuccessWithMessage(w, "status proyek berhasil diperbarui", nil)
}

func (h *UserProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id"), middlewares.GetUserID(r)); err != nil {
		if errors.Is(err, repositories.ErrUserProjectNotFound) {
			response.RespondNotFound(w, err.Error())
			return
		}
		response.RespondInternalServerError(w, err.Error())
		return
	}
	response.RespondSuccessWithMessage(w, "proyek berhasil diarsipkan", nil)
}
