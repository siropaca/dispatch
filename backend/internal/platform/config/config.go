package config

import (
	"fmt"
	"os"
)

// Config はアプリ全体のランタイム設定。環境変数から読み込む。
type Config struct {
	Env         string // "development" | "production"
	Port        string // HTTP リッスンポート
	DatabaseURL string // Postgres 接続文字列
}

// Load は環境変数から Config を構築する。必須値が無ければエラーを返す。
func Load() (Config, error) {
	cfg := Config{
		Env:         getenv("APP_ENV", "development"),
		Port:        getenv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("config: DATABASE_URL is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
