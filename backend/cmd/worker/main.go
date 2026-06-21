package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"

	"github.com/siropaca/dispatch/backend/internal/platform/telemetry"
	reportingv1 "github.com/siropaca/dispatch/backend/internal/proto/dispatch/reporting/v1"
	"github.com/siropaca/dispatch/backend/internal/proto/dispatch/reporting/v1/reportingv1connect"
)

// reportingServer は内部 Connect の ReportingService スタブ。
// Cloud Tasks から push されるジョブの受け口で、実処理は Phase 1 で実装する。
type reportingServer struct{}

func (reportingServer) RunReporting(
	_ context.Context,
	_ *connect.Request[reportingv1.RunReportingRequest],
) (*connect.Response[reportingv1.RunReportingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("reporting: not implemented yet"))
}

func main() {
	if err := run(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run() error {
	env := getenv("APP_ENV", "development")
	port := getenv("PORT", "8081")
	logger := telemetry.NewLogger(env)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	// 内部ジョブ契約(Connect)。Connect プロトコルは HTTP/1.1 で動く。
	path, handler := reportingv1connect.NewReportingServiceHandler(reportingServer{})
	mux.Handle(path, handler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("worker listening", "addr", srv.Addr, "env", env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
