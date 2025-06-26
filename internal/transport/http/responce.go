package http

import (
	"encoding/json"
	"net/http"
)

// Responder определяет интерфейс для отправки HTTP-ответов в виде JSON
type Responder interface {
	JSON(code int, body interface{})
	Error(code int, msg string)
}

// структура реализует интерфейс Responder, используя http.ResponseWriter
type jsonResponder struct {
	w http.ResponseWriter
}

// ErrorResponse описывает тело ответа с сообщением об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewResponder создаёт новый JSON-Responder для заданного ResponseWriter
func NewResponder(w http.ResponseWriter) Responder {
	return &jsonResponder{w: w}
}

// JSON отправляет HTTP-ответ с указанным кодом состояния и данными в формате JSON
func (r *jsonResponder) JSON(code int, body interface{}) {
	r.w.Header().Set("Content-Type", "application/json")
	if body == nil {
		r.w.WriteHeader(code)
		return
	}
	r.w.WriteHeader(code)
	_ = json.NewEncoder(r.w).Encode(body)
}

// Error отправляет JSON-ответ с сообщением об ошибке и соответствующим кодом
func (r *jsonResponder) Error(code int, msg string) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(code)
	_ = json.NewEncoder(r.w).Encode(ErrorResponse{Error: msg})
}
