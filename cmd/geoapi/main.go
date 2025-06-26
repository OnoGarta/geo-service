package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth/v5"

	"github.com/your-org/geo-service-swagger/internal/address"
	"github.com/your-org/geo-service-swagger/internal/auth"
	"github.com/your-org/geo-service-swagger/internal/dadata"
	thttp "github.com/your-org/geo-service-swagger/internal/transport/http"
)

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
	// Конфигурация сервиса: ключи DaData, секрет JWT, порт и таймауты
	cfg := struct {
		APIKey    string
		Secret    string
		JWTSecret string
		Port      string
		Timeout   time.Duration
	}{
		APIKey:    os.Getenv("DADATA_API_KEY"),
		Secret:    os.Getenv("DADATA_SECRET_KEY"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      "8080",
		Timeout:   3 * time.Second,
	}
	if cfg.APIKey == "" || cfg.Secret == "" {
		log.Fatal("DADATA_API_KEY and DADATA_SECRET_KEY must be set")
	}

	// Инициализация зависимостей: клиент DaData и сервис адресов
	client := dadata.New(cfg.APIKey, cfg.Secret, cfg.Timeout)
	geoSvc := address.NewService(client)
	userStore := auth.NewStore() // хранилище пользователей (в памяти)

	// Настройка JWT-аутентификации (секрет из переменной окружения или значение по умолчанию)
	secret := cfg.JWTSecret
	if secret == "" {
		secret = "super-secret-key"
		log.Println("Warning: using default JWT secret")
	}
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	// Создание роутера API с контроллерами
	router := thttp.Router(geoSvc, userStore, tokenAuth)
	server := NewAPIServer(":"+cfg.Port, router)

	// Запуск сервера и ожидание завершения (graceful shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("geo-service listening on :%s", cfg.Port)
		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-stop // ждём сигнала завершения

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
