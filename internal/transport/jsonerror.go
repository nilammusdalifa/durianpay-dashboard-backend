package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CodeToStatus(code entity.Code) int {
	switch code {
	case entity.ErrorCodeBadRequest:
		return http.StatusBadRequest
	case entity.ErrorCodeNotFound:
		return http.StatusNotFound
	case entity.ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func WriteAppError(w http.ResponseWriter, appErr *entity.AppError) {
	status := CodeToStatus(appErr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{
		Code:    status,
		Message: appErr.Message,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	var aErr *entity.AppError
	if errors.As(err, &aErr) {
		WriteAppError(w, aErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    500,
		Message: "internal error",
	})
}
