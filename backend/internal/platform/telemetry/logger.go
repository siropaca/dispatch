package telemetry

import (
	"log/slog"
	"os"
)

// NewLogger は env に応じた slog.Logger を返す。
// production は JSON、それ以外は人間が読みやすい text。
func NewLogger(env string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var h slog.Handler
	if env == "production" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(h)
}
