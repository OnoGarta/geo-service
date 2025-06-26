package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/your-org/geo-service-swagger/internal/address"
	"github.com/your-org/geo-service-swagger/internal/auth"
)

// Router настраивает маршруты API: Swagger UI, аутентификация и адресный сервис (с защитой JWT)
func Router(svc address.Service, store *auth.MemoryStore, tokenAuth *jwtauth.JWTAuth) chi.Router {
	r := chi.NewRouter()
	// Инициализируем контроллеры
	authCtrl := NewAuthController(store, tokenAuth)
	addrCtrl := NewAddressController(svc)
	// Открытые эндпойнты (доступ без токена)
	r.Mount("/swagger", SwaggerRouter())
	r.Post("/api/register", authCtrl.Register)
	r.Post("/api/login", authCtrl.Login)
	// Защищённые JWT эндпойнты для работы с адресами
	r.Route("/api/address", func(pr chi.Router) {
		pr.Use(auth.Middleware(tokenAuth))
		pr.Post("/search", addrCtrl.Search)
		pr.Post("/geocode", addrCtrl.Geocode)
	})
	return r
}
