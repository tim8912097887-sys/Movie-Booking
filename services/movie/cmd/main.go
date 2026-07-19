package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/configs"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/db"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/interface/http"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/shutdown"
)

func main() {

	// Initialize the logger
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
	// Initialize the config
	envConfigs := configs.InitConfigs()
	shutdownManager := shutdown.NewShutdownManager(logger)

	// Initialize the database connection pool
	dbConfig := db.DbConfig{Logger: logger, Ctx: ctx, DbUrl: envConfigs.Db.Url, ShutdownManager: shutdownManager}
	db := db.NewDb(dbConfig)
	pool,err := db.InitDB()

	if err != nil {
		logger.Error("failed to initialize the database connection pool",slog.Any("error", err))
		os.Exit(1)
		return
	}
	// Initialize the API
	api := http.Api{Config: http.ApiConfig{Logger: logger, EnvConfigs: envConfigs}}
	api.Run(ctx, api.Mount(pool), 8*time.Second,shutdownManager)
}