package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/siropaca/dispatch/backend/internal/platform/httpapi"
)

// GetReadyz は readiness を返す。依存(DB)へ疎通できれば 200、できなければ 503。
// 内部エラー詳細はクライアントへ漏らさない(status のみ返す)。
func (s apiServer) GetReadyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := s.db.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(httpapi.Health{Status: "unavailable"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(httpapi.Health{Status: "ok"})
}
