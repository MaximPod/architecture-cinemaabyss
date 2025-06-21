// This file is safe to edit. Once it exists it will not be overwritten
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"events/internal/config"
	"events/internal/pubsub"
)

// Service - приложение-сервис
type Service struct {
	ctx  context.Context // the main context
	pool Pool            // пул раннеров сервиса объединенных в errgroup

	Routers *pubsub.Routers // Роутеры и Шины транспортной системы PubSub

}

// Pool -  пул раннеров сервиса объединенных в errgroup
type Pool struct {
	runners  []func(ctx context.Context) (err error) // список раннеров, которые будут запущены
	shutdown chan struct{}                           // канал сигнала остановки раннеров сервиса
}

// New инициализирует сущность: сервис приложения
func New(ctx context.Context, config *config.Config) (*Service, error) {
	runners := make([]func(ctx context.Context) (err error), 0)
	shutdown := make(chan struct{})

	// транспорт PubSub
	routers, err := pubsub.NewRouters(ctx, config.PubSub)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewRouters:%w", err)
	}

	s := &Service{
		ctx: ctx,
		pool: Pool{
			runners:  runners,
			shutdown: shutdown,
		},
		Routers: routers,
	}

	return s, nil
}

// PreServerShutdown - обработчик остановки приложения - вызывается из сервера
func (s *Service) PreServerShutdown() {
	slog.Debug("service start PreServerShutdown")
	// сигнал завершения работы пула раннеров
	s.pool.shutdown <- struct{}{}
}

// OnShutdown - обработчик остановки приложения - вызывается из сервера
func (s *Service) OnShutdown() {
	slog.Debug("service start shutdown")
	s.Routers.MessageRouter.Close()

	close(s.pool.shutdown)
}

// Exit - логирует результат остановки приложения
func (s *Service) Exit(err error) {
	var es OSSignalError

	switch {
	case err == nil:
		slog.Info("app is stopped")
	case errors.As(err, &es):
		slog.Info("app is stopped by signal", "signal", es.signal.String())
	case errors.Is(err, http.ErrServerClosed):
		slog.Info("app is stopped by http server")
	case errors.Is(err, ShutdownCallError):
		slog.Info("app is stopped by http server")
	default:
		slog.Error("app is stopped", "error", err)
	}
}

// PubSubInitHandlers - подключаем обработчики PubSub
func (s *Service) PubSubInitHandlers() error {

	// event interface handlers
	err := s.Routers.EventProcessor.AddHandlers(
		pubsub.NewEventHandler("MovieEventHandler", s.MovieEventHandler),
		pubsub.NewEventHandler("UserEventHandler", s.UserEventHandler),
		pubsub.NewEventHandler("PaymentEventHandler", s.PaymentEventHandler),
	)
	if err != nil {
		return fmt.Errorf("s.Routers.EventProcessor.AddHandlers: %v", err)
	}

	return nil
}
