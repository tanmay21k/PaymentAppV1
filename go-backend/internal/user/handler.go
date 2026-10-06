package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tanmay21k/internal/database"
	"github.com/tanmay21k/internal/helpers"
)

func SayHello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.New().Info(" /hello endpoint hit")
		w.Write([]byte("Hello World!"))
	})
}

var log *slog.Logger = helpers.New()

func SignUp(queries *database.Queries) http.Handler {
	service := newService(newRepository(queries))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.ValidateMethod(log, w, r, http.MethodPost) {
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		request, err := validateSignUp(log, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := service.SignUp(r.Context(), request); err != nil {
			var postgresErr *pgconn.PgError
			if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
				http.Error(w, "username already exists", http.StatusConflict)
				return
			}
			log.Error("could not create user", "error", err)
			http.Error(w, "could not create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"firstName": request.FirstName,
			"lastName":  request.LastName,
			"username":  request.Username,
		})
	})
}
