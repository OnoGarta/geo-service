package http

import (
	"encoding/json"
	"net/http"
)

// JSONResponder реализует Responder для ответа в формате JSON.
type JSONResponder struct {
	Status int
	Body   any
}

func (r JSONResponder) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	_ = json.NewEncoder(w).Encode(r.Body)
}

// ErrorResponder реализует Responder для ошибок.
type ErrorResponder struct {
	Status  int
	Message string
}

func (r ErrorResponder) Respond(w http.ResponseWriter) {
	http.Error(w, r.Message, r.Status)
}
