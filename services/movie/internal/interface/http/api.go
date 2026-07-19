package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/configs"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/db"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/shutdown"
)

type ApiConfig struct {
	Logger *slog.Logger
	EnvConfigs configs.Configs
}

type Api struct {
	Config ApiConfig
}

func (a *Api) Mount() http.Handler {
	r := chi.NewRouter()

	// Register the middleware
	r.Use(middleware.Timeout(5*time.Second))
	// Health check
	r.Get("/health",func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize the movie repository, usecase and handler
	movieRepository := db.NewMovieRepository()
	createMovieUseCase := application.NewCreateMovieUsecase(movieRepository)
	getMovieUseCase := application.NewGetMovieUsecase(movieRepository)
	getMoviesUseCase := application.NewGetMoviesUsecase(movieRepository)
	updateMovieUseCase := application.NewUpdateMovieUsecase(movieRepository)
	deleteMovieUseCase := application.NewDeleteMovieUsecase(movieRepository)
	handlerConfig := HandlerConfig{
		Logger:              a.Config.Logger,
		CreateMovieUseCase:  createMovieUseCase,
		GetMovieUseCase:     getMovieUseCase,
		GetMoviesUseCase:    getMoviesUseCase,
		UpdateMovieUseCase:  updateMovieUseCase,
		DeleteMovieUseCase:  deleteMovieUseCase,	
	}
	handler := NewHandler(handlerConfig)
	r.Route("/api/v1/movies", func(r chi.Router) {
		handler.RegisterRoutes(r)	
	})
	return r
}

func (a *Api) Run(ctx context.Context, h http.Handler, shutdownTimeout time.Duration,shutdownManager *shutdown.ShutdownManager) error {
	server := &http.Server{
		Addr:    a.Config.EnvConfigs.Api.Addr,
		Handler: h,
		ReadTimeout:       5 * time.Second,
        ReadHeaderTimeout: 2 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,
	}

	// Channel to notify when the server is initialized failure
	serverErrorCh := make(chan error, 1)
	// Start the server with goroutine
	go func() {
		a.Config.Logger.Info("starting server",slog.String("address", a.Config.EnvConfigs.Api.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Config.Logger.Error("failed to start server",slog.Any("error", err))
			serverErrorCh <- err
		}
	}()

	// Register the shutdown handler
	shutdownManager.Register(a.Shutdown(server, shutdownTimeout))

	select {
		case <-ctx.Done():
			a.Config.Logger.Info("shutting down the server",slog.String("reason", ctx.Err().Error()))
		case err := <-serverErrorCh:
			return err
	}

	// Start a graceful shutdown
	shutdownManager.Shutdown(shutdownTimeout)

	return nil

}

func (a *Api) Shutdown(server *http.Server, shutdownTimeout time.Duration) func(context.Context) error {
	
	return func(ctx context.Context) error {

		if err := server.Shutdown(ctx); err != nil {
			a.Config.Logger.Error("failed to shut down the server",slog.Any("error", err))
			if closeErr := server.Close(); closeErr != nil {
				a.Config.Logger.Error("failed to close the server",slog.Any("error", err))
				return errors.Join(err,closeErr)
			}
			return err
		}

		a.Config.Logger.Info("server shut down gracefully")
		return nil
	}
}