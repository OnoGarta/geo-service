package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/your-org/geo-service/internal/address"
)

func SearchHandler(svc address.Service) http.HandlerFunc {
	v := validator.New()
	return func(w http.ResponseWriter, r *http.Request) {
		var req address.SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := v.Struct(req); err != nil {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}
		res, err := svc.Search(r.Context(), req.Query)
		if err != nil {
			http.Error(w, "upstream error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(address.Response{Addresses: res})
	}
}

func GeocodeHandler(svc address.Service) http.HandlerFunc {
	v := validator.New()
	return func(w http.ResponseWriter, r *http.Request) {
		var req address.GeocodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := v.Struct(req); err != nil {
			http.Error(w, "lat and lng required", http.StatusBadRequest)
			return
		}
		res, err := svc.Geocode(r.Context(), req.Lat, req.Lng)
		if err != nil {
			http.Error(w, "upstream error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(address.Response{Addresses: res})
	}
}
