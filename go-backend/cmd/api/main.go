package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tanmay21k/PaymentAppV1/go-backend/internal/config"
	"github.com/tanmay21k/PaymentAppV1/go-backend/internal/database"
	"github.com/tanmay21k/PaymentAppV1/go-backend/internal/health"
	"github.com/tanmay21k/PaymentAppV1/go-backend/internal/user"
)

func main() {
	// load the config from .env; a missing file is fine when variables come from the environment
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("[ERROR]: could not read .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR]: invalid configuration: %v", err)
	}
	if os.Getenv("PORT") == "" {
		log.Println("[WARN]: Using default port")
	}

	// Create the pool; the pool stays open for the life of the process
	pool, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[ERROR]: failed to connect to db, error: %v", err)
	}
	defer pool.Close()
	log.Println("[INFO]: PostgreSQL connection pool created")

	// Check the database now, then every 5 minutes, in the background
	go func() {
		time.Sleep(30 * time.Second)
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := health.CheckDBHealth(ctx, pool); err != nil {
				log.Printf("[ERROR]: database health check failed: %v", err)
				// intentional design to fail when dbConection is broke
				os.Exit(1)
			} else {
				log.Println("[INFO]: database health check passed")
			}
			cancel()
			<-ticker.C
		}
	}()

	// Register the health route on our own mux instead of the default one
	mux := http.NewServeMux()
	mux.Handle("GET /health", health.NewHealthClient(pool))
	mux.Handle("POST /api/v1/user/signup", user.SignUp())

	// Stop when the process gets Ctrl+C (SIGINT) or SIGTERM
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{Addr: ":" + cfg.Port, Handler: mux}
	go func() {
		log.Printf("[INFO]: listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[ERROR]: server failed: %v", err)
		}
	}()

	<-sigCtx.Done()
	log.Println("[INFO]: shutting down")

	// Give in-flight requests up to 10 seconds to finish
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR]: shutdown did not finish cleanly: %v", err)
	}
	log.Println("[INFO]: shutdown complete")
}
