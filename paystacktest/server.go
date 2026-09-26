package paystacktest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// Server is a minimal mock Paystack API server for tests, backed by
// httptest.Server. Register canned responses with Respond, then point a
// paystack.Client at the mock with paystack.WithBaseURL(server.URL).
type Server struct {
	*httptest.Server
	mux *http.ServeMux
}

// NewServer starts a mock server with no routes registered. Register
// routes with Respond before making requests against it, and call Close
// when the test finishes (typically via t.Cleanup).
func NewServer() *Server {
	mux := http.NewServeMux()
	s := &Server{mux: mux}
	s.Server = httptest.NewServer(mux)
	return s
}

// Respond registers a canned JSON response for method and path, wrapping
// data in the standard Paystack {status, message, data} response envelope.
func (s *Server) Respond(method, path string, statusCode int, data any) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  statusCode >= 200 && statusCode < 300,
			"message": http.StatusText(statusCode),
			"data":    data,
		})
	})
}
