package http

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func errBadRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func errInternal(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}
