package handlers

import (
	"encoding/json"
	"net/http"

	"go-produk/internal/dto"
	"go-produk/internal/middlewares"
	"go-produk/internal/services"
	"go-produk/pkg/response"
)

// UserHandler menangani request HTTP untuk user management
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler membuat instance UserHandler baru
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile menangani GET /api/v1/users/me
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	user, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccess(w, user)
}

// UpdateProfile menangani PUT /api/v1/users/me
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request tidak valid")
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	response.RespondSuccessWithMessage(w, "profil berhasil diperbarui", user)
}
