package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/siropaca/dispatch/backend/internal/platform/httpapi"
)

// Pinger は readiness で確認する依存(DB 等)の疎通チェック。
// *pgxpool.Pool が構造的にこれを満たす。
type Pinger interface {
	Ping(ctx context.Context) error
}

// apiServer は openapi.yaml から生成した httpapi.ServerInterface を実装する。
type apiServer struct {
	db Pinger
}

// newAPIServer は依存を注入して apiServer を構築する。
func newAPIServer(db Pinger) apiServer {
	return apiServer{db: db}
}

// GetHealthz は liveness を返す(依存に関係なく常に 200)。
func (s apiServer) GetHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(httpapi.Health{Status: "ok"})
}
