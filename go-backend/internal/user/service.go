package user

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tanmay21k/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type signUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func validateSignUp(log *slog.Logger, r *http.Request) (signUpRequest, error) {
	var request signUpRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		log.Error("error reading body", "error", err)
		return signUpRequest{}, errors.New("invalid request body")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return signUpRequest{}, errors.New("request body must contain a single JSON object")
	}

	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Username = strings.TrimSpace(request.Username)
	if request.FirstName == "" || request.LastName == "" || request.Username == "" || strings.TrimSpace(request.Password) == "" {
		log.Error("bad request")
		return signUpRequest{}, errors.New("all fields are required")
	}
	if len(request.Password) > 72 {
		return signUpRequest{}, errors.New("password must be 72 bytes or fewer")
	}

	return request, nil
}

func newService(repo repository) *svc {
	return &svc{repo: repo}
}

func (s *svc) SignUp(ctx context.Context, request signUpRequest) (database.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, err
	}

	return s.repo.CreateUser(ctx, database.CreateUserParams{
		Firstname: pgtype.Text{String: request.FirstName, Valid: true},
		Lastname:  pgtype.Text{String: request.LastName, Valid: true},
		Username:  request.Username,
		Password:  pgtype.Text{String: string(passwordHash), Valid: true},
	})
}
