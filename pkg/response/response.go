package response

import (
	"encoding/json"
	"net/http"
)

// ResponseBody adalah representasi standar JSON response untuk frontend
type ResponseBody struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// RespondJSON menulis response JSON dengan status code tertentu
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// RespondSuccess merespon dengan sukses (200 OK)
func RespondSuccess(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, ResponseBody{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "sukses",
		Data:       data,
	})
}

// RespondSuccessWithMessage merespon dengan sukses (200 OK) beserta custom message
func RespondSuccessWithMessage(w http.ResponseWriter, message string, data interface{}) {
	RespondJSON(w, http.StatusOK, ResponseBody{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	})
}

// RespondCreated merespon dengan created (201 Created)
func RespondCreated(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusCreated, ResponseBody{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "berhasil dibuat",
		Data:       data,
	})
}

// RespondError merespon dengan format error standar
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ResponseBody{
		Success:    false,
		StatusCode: status,
		Message:    message,
		Data:       nil,
	})
}

// RespondBadRequest merespon 400 Bad Request
func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusBadRequest, message)
}

// RespondUnauthorized merespon 401 Unauthorized
func RespondUnauthorized(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusUnauthorized, message)
}

// RespondForbidden merespon 403 Forbidden
func RespondForbidden(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusForbidden, message)
}

// RespondNotFound merespon 404 Not Found
func RespondNotFound(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusNotFound, message)
}

// RespondInternalServerError merespon 500 Internal Server Error
func RespondInternalServerError(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusInternalServerError, message)
}
