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

	// openapi.yaml(spec-first)から生成した interface を /api 配下に結線する。
	// 公開は web(Caddy)が単一オリジンで /api/* を api へ振り分ける(ADR-0015)。
	httpapi.HandlerFromMuxWithBaseURL(newAPIServer(db), r, "/api")

	_ = logger // ルート別ロギング等で利用予定(今後拡張)
	return r
}
