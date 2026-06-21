package config

import "testing"

func TestLoad(t *testing.T) {
	t.Run("applies defaults", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://localhost/x")
		t.Setenv("APP_ENV", "")
		t.Setenv("PORT", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Env != "development" {
			t.Errorf("Env = %q, want development", cfg.Env)
		}
		if cfg.Port != "8080" {
			t.Errorf("Port = %q, want 8080", cfg.Port)
		}
	})

	t.Run("env overrides defaults", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://localhost/x")
		t.Setenv("APP_ENV", "production")
		t.Setenv("PORT", "9999")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Env != "production" || cfg.Port != "9999" {
			t.Errorf("got %+v", cfg)
		}
	})

	t.Run("missing DATABASE_URL is an error", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")
		if _, err := Load(); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
