package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/your-org/geo-service/internal/address"
)

func Router(svc address.Service) chi.Router {
	r := chi.NewRouter()
	r.Post("/api/address/search", SearchHandler(svc))
	r.Post("/api/address/geocode", GeocodeHandler(svc))

	r.Post("/address/search", SearchHandler(svc))
	r.Post("/address/geocode", GeocodeHandler(svc))
	return r
}
