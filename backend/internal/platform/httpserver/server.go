package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/siropaca/dispatch/backend/internal/platform/httpapi"
)

// NewRouter はミドルウェアと、openapi 契約から生成したルートを備えた http.Handler を構築する。
// db は readiness(/readyz)で疎通確認する依存。
func NewRouter(logger *slog.Logger, db Pinger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// openapi.yaml(spec-first)から生成した interface を実装に結線する。
	httpapi.HandlerFromMux(newAPIServer(db), r)

	_ = logger // ルート別ロギング等で利用予定(今後拡張)
	return r
}
