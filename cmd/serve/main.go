package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/app"
	"github.com/PhilippReinke/kino-berlin/pkg/delivery"
	"github.com/PhilippReinke/kino-berlin/pkg/domain"
	"github.com/PhilippReinke/kino-berlin/pkg/infra/provider"
	"github.com/PhilippReinke/kino-berlin/pkg/infra/storage"
)

func main() {
	cfg := parseFlags()

	storage := storage.NewMemory()
	babylon := provider.NewBabylon()
	yorck := provider.NewYorck()

	application := app.New(
		storage,
		[]domain.Provider{babylon, yorck},
		app.Config{
			SyncInterval: cfg.SyncInterval,
		},
	)

	if err := application.StartBackgroundSync(); err != nil {
		log.Fatalf("Failed to start background sync: %v", err)
	}

	handler, err := delivery.NewHandler(
		application,
		cfg.TemplateDir,
		cfg.StaticDir,
	)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	if err := runServer(cfg.Addr, handler, application); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runServer(addr string, handler *delivery.Handler, application *app.App) error {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on http://%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	application.StopBackgroundSync()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}
