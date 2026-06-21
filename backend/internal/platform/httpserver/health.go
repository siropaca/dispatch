package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/siropaca/dispatch/backend/internal/platform/httpapi"
)

// apiServer は openapi.yaml から生成した httpapi.ServerInterface を実装する。
type apiServer struct{}

// GetHealthz は liveness を返す(常に 200)。
func (apiServer) GetHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(httpapi.Health{Status: "ok"})
}
