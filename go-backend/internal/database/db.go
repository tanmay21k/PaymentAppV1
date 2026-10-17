package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB creates the connection pool and returns it to the caller, who
// is responsible for closing it. pgxpool.New only parses the URL and builds
// the pool; it does not dial the database. Reachability is checked by the
// health check, not here.
func ConnectDB(dbUrl string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
