package error_handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func JsonError(w http.ResponseWriter, message string, code int, details ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var det interface{}
	if len(details) > 0 {
		det = details[0]
	}

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    http.StatusText(code),
		Details: det,
	})
}
