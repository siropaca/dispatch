package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	apiServer{}.GetHealthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if rec.Body.String() == "" {
		t.Error("body is empty")
	}
}
