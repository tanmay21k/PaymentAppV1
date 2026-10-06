package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/tanmay21k/internal/database"
	logger "github.com/tanmay21k/internal/helpers"
	"github.com/tanmay21k/internal/user"
	"github.com/tanmay21k/internal/utils"
)

func main() {
	log := logger.New()
	log.Info("starting application")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := utils.ConnectDB(ctx)
	cancel()
	if err != nil {
		log.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	log.Info("database connection established")

	queries := database.New(pool)
	mux := http.NewServeMux()
	mux.Handle("/hello", user.SayHello())
	mux.Handle("/signup", user.SignUp(queries))

	log.Info("server started", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("server has stopped", "error", err)
		os.Exit(1)
	}
}
