package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHealthReturns503WhenDatabaseUnreachable(t *testing.T) {
	// Nothing listens on port 1, so the ping fails without needing a database.
	pool, err := pgxpool.New(context.Background(), "postgres://user:pass@127.0.0.1:1/db?connect_timeout=1")
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	NewHealthClient(pool).ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "unavailable" || body.Database != "unreachable" {
		t.Fatalf("body = %+v, want status=unavailable database=unreachable", body)
	}
}

func TestHealthReturns503WhenDatabaseIsSlow(t *testing.T) {
	// Accept connections but never reply, so the ping hangs until the handler's timeout.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Drain and hold the connection open without writing anything.
			go io.Copy(io.Discard, conn)
		}
	}()

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://user:pass@%s/db", ln.Addr()))
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	start := time.Now()
	NewHealthClient(pool).ServeHTTP(rec, req)
	elapsed := time.Since(start)

	// The handler must answer near pingTimeout, not wait for the database.
	if elapsed > pingTimeout+time.Second {
		t.Fatalf("handler took %v, want about %v", elapsed, pingTimeout)
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "unavailable" || body.Database != "unreachable" {
		t.Fatalf("body = %+v, want status=unavailable database=unreachable", body)
	}
}
