package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CheckDBHealth pings the database through the shared pool and returns the
// result. The caller decides what to do with a failure.
func CheckDBHealth(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}
