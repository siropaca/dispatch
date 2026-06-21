package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/siropaca/dispatch/backend/internal/platform/httpapi"
)

// NewRouter はミドルウェアと、openapi 契約から生成したルートを備えた http.Handler を構築する。
func NewRouter(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// openapi.yaml(spec-first)から生成した interface を実装に結線する。
	httpapi.HandlerFromMux(apiServer{}, r)

	_ = logger // ルート別ロギング等で利用予定(今後拡張)
	return r
}
