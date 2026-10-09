package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ---- Service ----

func newService(repo repository) *svc {
	return &svc{repo: repo}
}

func (s *svc) SignUp(ctx context.Context, request signUpRequest) (database.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, database.CreateUserParams{
		Firstname: pgtype.Text{String: request.FirstName, Valid: true},
		Lastname:  pgtype.Text{String: request.LastName, Valid: true},
		Username:  request.Username,
		Password:  pgtype.Text{String: string(passwordHash), Valid: true},
	})
	if err != nil {
		return database.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

var errInvalidCredentials = errors.New("invalid credentials")

func (s *svc) SignIn(ctx context.Context, request signInRequest) error {
	user, err := s.repo.FetchUser(ctx, request.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errInvalidCredentials
		}
		return fmt.Errorf("fetch user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password.String),
		[]byte(request.Password),
	)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return errInvalidCredentials
	}
	if err != nil {
		return fmt.Errorf("compare password: %w", err)
	}
	return nil
}
