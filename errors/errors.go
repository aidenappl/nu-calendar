package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func SendError(w http.ResponseWriter, err string, message string, status int) {
	errResp := ErrorResponse{
		Error:   err,
		Message: http.StatusText(status),
		Status:  status,
	}

	if len(message) != 0 {
		errResp.Message = message
	} else {
		errResp.Message = http.StatusText(status)
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}
