package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/your-org/geo-service/internal/address"
	"github.com/your-org/geo-service/internal/auth"
	"github.com/your-org/geo-service/internal/dadata"
	thttp "github.com/your-org/geo-service/internal/transport/http"
)

var tokenAuth = jwtauth.New("HS256", []byte("super-secret-key"), nil)

func main() {
	cfg := struct {
		APIKey  string
		Secret  string
		Port    string
		Timeout time.Duration
	}{
		APIKey:  os.Getenv("DADATA_API_KEY"),
		Secret:  os.Getenv("DADATA_SECRET_KEY"),
		Port:    "8080",
		Timeout: 3 * time.Second,
	}
	if cfg.APIKey == "" || cfg.Secret == "" {
		log.Fatal("DADATA_API_KEY and DADATA_SECRET_KEY must be set")
	}

	client := dadata.New(cfg.APIKey, cfg.Secret, cfg.Timeout)
	geoSvc := address.NewService(client)
	userStore := auth.NewStore() // in-memory store

	r := chi.NewRouter()

	// --- Открытые эндпойнты ---
	r.Mount("/swagger", thttp.SwaggerRouter())
	r.Post("/api/register", auth.Register(userStore))
	r.Post("/api/login", auth.Login(userStore, tokenAuth))

	// --- Группа защищённых адресных эндпойнтов ---
	r.Route("/api/address", func(protected chi.Router) {
		protected.Use(auth.Middleware(tokenAuth))
		protected.Post("/search", thttp.SearchHandler(geoSvc))
		protected.Post("/geocode", thttp.GeocodeHandler(geoSvc))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("geo-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("bye")
}
