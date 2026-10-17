package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		port        string
		databaseURL string
		wantPort    string
		wantErr     bool
		errMentions string
	}{
		{name: "uses default port when PORT unset", port: "", databaseURL: "postgres://x", wantPort: "4500"},
		{name: "uses PORT when set", port: "8080", databaseURL: "postgres://x", wantPort: "8080"},
		{name: "accepts lowest valid port", port: "1", databaseURL: "postgres://x", wantPort: "1"},
		{name: "accepts highest valid port", port: "65535", databaseURL: "postgres://x", wantPort: "65535"},
		{name: "rejects non-numeric port", port: "abc", databaseURL: "postgres://x", wantErr: true, errMentions: "PORT"},
		{name: "rejects port 0", port: "0", databaseURL: "postgres://x", wantErr: true, errMentions: "PORT"},
		{name: "rejects port above 65535", port: "70000", databaseURL: "postgres://x", wantErr: true, errMentions: "PORT"},
		{name: "rejects missing DATABASE_URL", port: "", databaseURL: "", wantErr: true, errMentions: "DATABASE_URL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)
			t.Setenv("DATABASE_URL", tt.databaseURL)

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Load() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.errMentions) {
					t.Fatalf("error %q does not mention %q", err, tt.errMentions)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if cfg.Port != tt.wantPort {
				t.Fatalf("Port = %q, want %q", cfg.Port, tt.wantPort)
			}
			if cfg.DatabaseURL != tt.databaseURL {
				t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, tt.databaseURL)
			}
		})
	}
}
