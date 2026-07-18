package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/configs"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/interface/http"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/shutdown"
)

func main() {

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	envConfigs := configs.InitConfigs()
	shutdownManager := shutdown.NewShutdownManager(logger)

	api := http.Api{Config: http.ApiConfig{Logger: logger, EnvConfigs: envConfigs}}
	api.Run(ctx, api.Mount(), 8*time.Second,shutdownManager)
}