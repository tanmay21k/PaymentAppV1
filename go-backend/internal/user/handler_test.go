package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty fails", "", false},
		{"whitespace passes", "   ", true},
		{"single letter passes", "A", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidName(tt.input); got != tt.want {
				t.Errorf("isValidName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"5 characters fails", "abcde", false},
		{"6 characters passes", "abcdef", true},
		{"6 characters with a trailing space passes", "abcde ", true},
		{"6 non-ASCII characters passes", "pässwö", true},
		{"5 non-ASCII characters fails (10 bytes)", "ééééé", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPassword(tt.input); got != tt.want {
				t.Errorf("isValidPassword(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	longLocal := strings.Repeat("a", 64)
	tooLongTotal := longLocal + "@" + strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + ".com"

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"simple address passes", "a@example.com", true},
		{"allowed local symbols and subdomains pass", "first.last+tag_1%x-y@sub.example.co.uk", true},
		{"two-letter TLD passes", "a@b.co", true},
		{"uppercase passes", "A1@Example.COM", true},
		{"64-character local part passes", longLocal + "@example.com", true},

		{"empty fails", "", false},
		{"no @ fails", "not-an-email", false},
		{"two @ fails", "a@b@example.com", false},
		{"space fails", "a b@example.com", false},
		{"disallowed character fails", "a!b@example.com", false},
		{"leading dot in local part fails", ".a@example.com", false},
		{"trailing dot in local part fails", "a.@example.com", false},
		{"doubled dot in local part fails", "a..b@example.com", false},
		{"65-character local part fails", strings.Repeat("a", 65) + "@example.com", false},
		{"total length over 254 fails", tooLongTotal, false},
		{"empty local part fails", "@example.com", false},
		{"domain without a dot fails", "a@localhost", false},
		{"label starting with hyphen fails", "a@-example.com", false},
		{"label ending with hyphen fails", "a@example-.com", false},
		{"empty label fails", "a@example..com", false},
		{"trailing dot fails", "a@example.com.", false},
		{"underscore in domain fails", "a@ex_ample.com", false},
		{"one-letter TLD fails", "a@example.c", false},
		{"TLD with a digit fails", "a@example.c0m", false},
		{"numeric TLD fails", "a@example.123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUsername(tt.input); got != tt.want {
				t.Errorf("isValidUsername(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSignUpErrorBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantMsg string
	}{
		{"missing field", `{"lastname":"B","username":"a@example.com","password":"secret1"}`, "Invalid inputs"},
		{"5-character password", `{"firstname":"A","lastname":"B","username":"a@example.com","password":"abcde"}`, "Invalid inputs"},
		{"bad email", `{"firstname":"A","lastname":"B","username":"not-an-email","password":"secret1"}`, "Invalid inputs"},
		{"malformed JSON", `{"firstname":`, "invalid request body"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/signup", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			SignUp()(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("body is not JSON: %v (body %q)", err, rec.Body.String())
			}
			if got["msg"] != tt.wantMsg {
				t.Errorf("msg = %q, want %q", got["msg"], tt.wantMsg)
			}
		})
	}
}

func TestSignUpNameCheck(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "empty lastname returns 400",
			body:       `{"firstname":"A","lastname":"","username":"a@example.com","password":"secret1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty firstname returns 400",
			body:       `{"firstname":"","lastname":"B","username":"a@example.com","password":"secret1"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "whitespace lastname is not rejected by the name check",
			body:       `{"firstname":"A","lastname":"   ","username":"a@example.com","password":"secret1"}`,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/signup", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			SignUp()(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
