package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/geo-service-swagger/internal/address"
)

// Мокаем сервис
type mockService struct{}

func (m *mockService) Search(ctx context.Context, q string) ([]*address.Address, error) {
	return []*address.Address{{City: "Москва"}}, nil
}
func (m *mockService) Geocode(ctx context.Context, lat, lng string) ([]*address.Address, error) {
	return []*address.Address{{City: "Сочи"}}, nil
}

type mockServiceError struct{}

func (m *mockServiceError) Search(ctx context.Context, q string) ([]*address.Address, error) {
	return nil, fmt.Errorf("err")
}
func (m *mockServiceError) Geocode(ctx context.Context, lat, lng string) ([]*address.Address, error) {
	return nil, fmt.Errorf("err")
}

func TestSearchHandler_OK(t *testing.T) {
	handler := SearchHandler(&mockService{})
	reqBody, _ := json.Marshal(address.SearchRequest{Query: "Москва"})
	req := httptest.NewRequest("POST", "/api/address/search", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Москва")) {
		t.Fatalf("want response to contain 'Москва', got %s", w.Body.String())
	}
}

func TestGeocodeHandler_OK(t *testing.T) {
	handler := GeocodeHandler(&mockService{})
	reqBody, _ := json.Marshal(address.GeocodeRequest{Lat: "43.6", Lng: "39.7"})
	req := httptest.NewRequest("POST", "/api/address/geocode", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Сочи")) {
		t.Fatalf("want response to contain 'Сочи', got %s", w.Body.String())
	}
}

func TestSearchHandler_BadRequest(t *testing.T) {
	handler := SearchHandler(&mockService{})

	req := httptest.NewRequest("POST", "/api/address/search", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	req2 := httptest.NewRequest("POST", "/api/address/search", bytes.NewBufferString(`{"query":""}`))
	w2 := httptest.NewRecorder()
	handler(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w2.Code)
	}
}

func TestGeocodeHandler_BadRequest(t *testing.T) {
	handler := GeocodeHandler(&mockService{})

	req := httptest.NewRequest("POST", "/api/address/geocode", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	req2 := httptest.NewRequest("POST", "/api/address/geocode", bytes.NewBufferString(`{"lat":""}`))
	w2 := httptest.NewRecorder()
	handler(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w2.Code)
	}
}

func TestSearchHandler_InternalError(t *testing.T) {
	handler := SearchHandler(&mockServiceError{})
	body, _ := json.Marshal(address.SearchRequest{Query: "err"})
	req := httptest.NewRequest("POST", "/api/address/search", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGeocodeHandler_InternalError(t *testing.T) {
	handler := GeocodeHandler(&mockServiceError{})
	body, _ := json.Marshal(address.GeocodeRequest{Lat: "43.6", Lng: "39.7"})
	req := httptest.NewRequest("POST", "/api/address/geocode", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
