package app

import (
	"context"
	"log/slog"
)

// MovieEventHandler - обработчик события MovieEvent
func (s *Service) MovieEventHandler(ctx context.Context, msg *MovieEvent) error {
	if msg == nil {
		slog.Debug("App: MovieEventHandler", "msg", "msg is nil")

		return nil
	}

	slog.Info("App: MovieEventHandler", "msg", *msg)

	return nil
}

// UserEventHandler - обработчик события UserEvent
func (s *Service) UserEventHandler(ctx context.Context, msg *UserEvent) error {
	if msg == nil {
		slog.Debug("App: UserEventHandler", "msg", "msg is nil")

		return nil
	}

	slog.Info("App: UserEventHandler", "msg", *msg)

	return nil
}

// PaymentEventHandler - обработчик события PaymentEvent
func (s *Service) PaymentEventHandler(ctx context.Context, msg *PaymentEvent) error {
	if msg == nil {
		slog.Debug("App: PaymentEventHandler", "msg", "msg is nil")

		return nil
	}

	slog.Info("App: PaymentEventHandler", "msg", *msg)

	return nil
}
