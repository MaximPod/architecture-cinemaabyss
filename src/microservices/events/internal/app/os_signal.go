// This file is safe to edit. Once it exists it will not be overwritten
package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// OSSignalError - сигнал ОС в виде ошибки для того, чтобы использовать SignalNotifyRun в errgroup
type OSSignalError struct {
	signal os.Signal
}

// Error - реализация интерфейса error
func (e OSSignalError) Error() string {
	return e.signal.String()
}

// SignalNotifyRun - раннер слушает сигналы ОС (SIGINT, SIGTERM) и возвращает полученный сигнал в виде ошибки
func (s *Service) SignalNotifyRun(ctx context.Context) error {
	slog.Debug("App: SignalNotifyRun is run")

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		slog.Debug("App: SignalNotifyRun is finished")
		return ctx.Err()

	case sgn := <-c:
		return OSSignalError{
			signal: sgn,
		}
	}
}
