package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/your-org/geo-service-swagger/internal/address"
	"github.com/your-org/geo-service-swagger/internal/auth"
	"github.com/your-org/geo-service-swagger/internal/dadata"
	thttp "github.com/your-org/geo-service-swagger/internal/transport/http"
)

var tokenAuth = jwtauth.New("HS256", []byte("super-secret-key"), nil)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(addr string, handler http.Handler) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *APIServer) Serve() error {
	return s.httpServer.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

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

	server := NewAPIServer(":"+cfg.Port, r)

	// --- graceful shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("geo-service listening on :%s", cfg.Port)
		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-stop // Ожидание сигнала

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
