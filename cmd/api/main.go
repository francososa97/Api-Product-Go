// Command api es el punto de entrada de la Product API. Se encarga del armado
// de dependencias (composition root) y del ciclo de vida del servidor HTTP,
// incluyendo el apagado ordenado (graceful shutdown).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/francososa97/product-api/internal/application"
	"github.com/francososa97/product-api/internal/domain"
	"github.com/francososa97/product-api/internal/infrastructure/config"
	"github.com/francososa97/product-api/internal/infrastructure/persistence/memory"
	mongorepo "github.com/francososa97/product-api/internal/infrastructure/persistence/mongo"
	httpserver "github.com/francososa97/product-api/internal/presentation/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error fatal: %v", err)
	}
}

// run concentra el arranque y devuelve error en lugar de hacer log.Fatal en
// cada paso, lo que facilita el testeo y un apagado limpio.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Selección del repositorio según la configuración (in-memory o MongoDB).
	repo, cleanup, err := buildRepository(cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	service := application.NewProductService(repo)
	router := httpserver.NewRouter(service)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.RequestTimeout,
		WriteTimeout: cfg.RequestTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// Se escucha en una goroutine para poder esperar señales de apagado.
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("servidor escuchando en :%s (driver=%s)", cfg.Port, cfg.Driver)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Espera a un error del servidor o a una señal de interrupción del SO.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Println("apagando servidor...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("servidor apagado correctamente")
	return nil
}

// buildRepository construye la implementación de persistencia según la config y
// devuelve también una función de limpieza para liberar recursos al apagar.
func buildRepository(cfg *config.Config) (domain.ProductRepository, func(), error) {
	switch cfg.Driver {
	case config.DriverMongo:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		repo, err := mongorepo.NewProductRepository(ctx, cfg.MongoURI, cfg.MongoDB, cfg.MongoCollection)
		if err != nil {
			return nil, nil, err
		}

		cleanup := func() {
			closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer closeCancel()
			if err := repo.Close(closeCtx); err != nil {
				log.Printf("error al cerrar la conexión con MongoDB: %v", err)
			}
		}
		return repo, cleanup, nil

	default: // config.DriverMemory
		return memory.NewProductRepository(), func() {}, nil
	}
}
