package health

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// pingTimeout is the longest /health waits for the database before answering 503.
const pingTimeout = 2 * time.Second

// healthResponse fields must be exported so encoding/json can see them.
type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// HealthClient serves GET /health. It is an http.Handler, so the server
// calls ServeHTTP once per request with that request's writer and request.
type HealthClient struct {
	pool *pgxpool.Pool
}

// NewHealthClient builds the handler with the pool it should check.
func NewHealthClient(pool *pgxpool.Pool) HealthClient {
	return HealthClient{pool: pool}
}

func (h HealthClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO]: /health endpoint")
	// Bound the ping so a hung database cannot hold the request open
	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()
	if err := CheckDBHealth(ctx, h.pool); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		err := json.NewEncoder(w).Encode(healthResponse{
			Status:   "unavailable",
			Database: "unreachable",
		})
		if err != nil {
			log.Printf("[ERROR]: failed to write health response, err: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(healthResponse{
		Status:   "ok",
		Database: "ok",
	})
	if err != nil {
		log.Printf("[ERROR]: failed to write health response, err: %v", err)
	}
}
