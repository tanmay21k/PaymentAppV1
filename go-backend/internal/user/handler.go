package user

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tanmay21k/internal/database"
	"github.com/tanmay21k/internal/helpers"
)

func SayHello(log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("hello endpoint hit", "method", r.Method, "path", r.URL.Path)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			log.Error("could not write hello response", "error", err)
		}
	})
}

func SignUp(log *slog.Logger, queries *database.Queries) http.Handler {
	service := newService(newRepository(queries))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.ValidateMethod(log, w, r, http.MethodPost) {
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var request signUpRequest
		if err := decodeJSON(r, &request); err != nil {
			writeRequestBodyError(w, err)
			return
		}
		request, err := validateSignUp(request)
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
		if err := json.NewEncoder(w).Encode(map[string]string{
			"firstName": request.FirstName,
			"lastName":  request.LastName,
			"username":  request.Username,
		}); err != nil {
			log.Error("could not write signup response", "error", err)
		}
	})
}

func SignIn(log *slog.Logger, queries *database.Queries) http.Handler {
	service := newService(newRepository(queries))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.ValidateMethod(log, w, r, http.MethodPost) {
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		var request signInRequest
		if err := decodeJSON(r, &request); err != nil {
			writeRequestBodyError(w, err)
			return
		}
		request.Username = strings.TrimSpace(request.Username)
		if request.Username == "" || strings.TrimSpace(request.Password) == "" {
			http.Error(w, "username and password are required", http.StatusBadRequest)
			return
		}

		if err := service.SignIn(r.Context(), request); err != nil {
			if errors.Is(err, errInvalidCredentials) {
				http.Error(w, "invalid username or password", http.StatusUnauthorized)
				return
			}
			log.Error("could not sign in user", "error", err)
			http.Error(w, "could not sign in", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		// TODO: Return a JWT token in the response.
	})
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err != nil {
			return err
		}
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

func writeRequestBodyError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}
	http.Error(w, "invalid request body", http.StatusBadRequest)
}

func validateSignUp(request signUpRequest) (signUpRequest, error) {
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Username = strings.TrimSpace(request.Username)
	if request.FirstName == "" || request.LastName == "" || request.Username == "" || strings.TrimSpace(request.Password) == "" {
		return signUpRequest{}, errors.New("all fields are required")
	}
	if len(request.Password) > 72 {
		return signUpRequest{}, errors.New("password must be 72 bytes or fewer")
	}
	return request, nil
}
