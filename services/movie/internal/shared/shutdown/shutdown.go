package shutdown

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type ShutdownManager struct {
	logger *slog.Logger
    mu       sync.Mutex
	shutdownFunc []func(context.Context) error
}

func NewShutdownManager(logger *slog.Logger) *ShutdownManager {
	return &ShutdownManager{
		shutdownFunc: make([]func(context.Context) error, 0),
	    logger: logger,
		mu : sync.Mutex{},
	}
}

func (s *ShutdownManager) Register(f func(context.Context) error) {
	s.mu.Lock()	
	defer s.mu.Unlock()
	s.shutdownFunc = append(s.shutdownFunc, f)
}

func (s *ShutdownManager) Shutdown(timeout time.Duration) error {
	s.mu.Lock()
	ctx, cancel := context.WithTimeout(context.Background(),timeout)
    defer cancel()

	funcs := make([]func(context.Context) error, len(s.shutdownFunc))
    copy(funcs, s.shutdownFunc)
	s.mu.Unlock()

	slices.Reverse(funcs)

	for _, f := range funcs {
		if err := f(ctx); err != nil {
			s.logger.Error("failed to shutdown",slog.Any("error", err))
			return err
		}
	}

	return nil
}