package httpserver

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakePinger は DB 疎通の成否を決め打ちする Pinger のテストダブル。
type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestGetReadyz(t *testing.T) {
	t.Run("returns 200 when dependency is reachable", func(t *testing.T) {
		srv := newAPIServer(fakePinger{err: nil})
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		srv.GetReadyz(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if body := rec.Body.String(); !strings.Contains(body, "ok") {
			t.Errorf("body = %q, want to contain %q", body, "ok")
		}
	})

	t.Run("returns 503 when dependency is unreachable", func(t *testing.T) {
		srv := newAPIServer(fakePinger{err: errors.New("connection refused")})
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		srv.GetReadyz(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
		if body := rec.Body.String(); !strings.Contains(body, "unavailable") {
			t.Errorf("body = %q, want to contain %q", body, "unavailable")
		}
		// 内部エラー詳細をクライアントへ漏らさない。
		if body := rec.Body.String(); strings.Contains(body, "connection refused") {
			t.Errorf("body leaks internal error: %q", body)
		}
	})
}

// TestRouter_Readyz は openapi 経由のルート結線(/readyz が NewRouter で配信される)を検証する。
func TestRouter_Readyz(t *testing.T) {
	ts := httptest.NewServer(NewRouter(discardLogger(), fakePinger{err: nil}))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("GET /readyz: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
}
