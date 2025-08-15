package errors

import (
	"encoding/json"
	"net/http"
)

func WriteHTTP(w http.ResponseWriter, appErr *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	_ = json.NewEncoder(w).Encode(appErr)
}
