package user

import (
	"context"

	"github.com/tanmay21k/internal/database"
)

type repository interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
}

type sqlcRepository struct {
	queries *database.Queries
}

func newRepository(queries *database.Queries) repository {
	return sqlcRepository{queries: queries}
}

func (r sqlcRepository) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

type svc struct {
	repo repository
}
