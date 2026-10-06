package helpers

import (
	"log/slog"
	"net/http"
)

func ValidateMethod(log *slog.Logger, w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		log.Info("method not allowed", "method", r.Method)

		w.Header().Set("Allow", expectedMethod)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return false
	}

	return true

}
