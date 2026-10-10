package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tanmay21k/internal/auth"
	"github.com/tanmay21k/internal/database"
	logger "github.com/tanmay21k/internal/helpers"
	"github.com/tanmay21k/internal/user"
	"github.com/tanmay21k/internal/utils"
)

func main() {
	log := logger.New()
	log.Info("starting application")
	if err := godotenv.Load(); err != nil {
		log.Info(".env file not found; using process environment")
	}

	tokens, err := auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Error("invalid JWT configuration", "error", err)
		os.Exit(1)
	}

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
	mux.Handle("/hello", user.SayHello(log))
	signUpHandler := user.SignUp(log, queries, tokens)
	signInHandler := user.SignIn(log, queries, tokens)
	mux.Handle("/signup", signUpHandler)
	mux.Handle("/signin", signInHandler)
	mux.Handle("/api/v1/user/signup", signUpHandler)
	mux.Handle("/api/v1/user/signin", signInHandler)

	log.Info("server started", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("server has stopped", "error", err)
		os.Exit(1)
	}
}
