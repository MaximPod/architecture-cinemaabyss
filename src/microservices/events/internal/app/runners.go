// This file is safe to edit. Once it exists it will not be overwritten
package app

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

var (
	// ShutdownCallError - псевдоошибка завершения работы сервиса
	ShutdownCallError = errors.New("Internal shutdown call")
)

// AddRunners добавляет раннеры в сервис.
// Раннер - блокирующая функция. См. функцию StartRunners.
func (s *Service) AddRunners(rs ...func(ctx context.Context) (err error)) {
	s.pool.runners = append(s.pool.runners, rs...)
}

// StartRunners - запускает раннеры конкурентно
// Раннер - блокирующая функция.
// @param ctx - контекст errgroup
// @return error - ошибка функции раннера или ctx.Err() если контекст завершен
// возврат ошибки из любого раннера приводит к завершению работы всех раннеров сервиса и самого приложения
func (s *Service) StartRunners() error {
	group, ctx := errgroup.WithContext(s.ctx)

	for _, fn := range s.pool.runners {
		internalFn := fn

		group.Go(func() error {
			return internalFn(ctx)
		})
	}

	return group.Wait()
}

// HandleShutdown - раннер, при получении сигнала в канале shutdown
// генерирует ошибку ShutdownCallError в errgroup для завершения пула раннеров
func (s *Service) HandleShutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.pool.shutdown:
		slog.Debug("App: s.pool.shutdown signal recieved")
		return ShutdownCallError
	}
}
