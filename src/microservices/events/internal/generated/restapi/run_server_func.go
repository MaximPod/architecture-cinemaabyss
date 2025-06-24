// This file is safe to edit. Once it exists it will not be overwritten
package restapi

import (
	"context"
)

// RunServerFunc - раннер для запуска web сервера
// при закрытии контекста выполняет graceful shutdown
func (s *Server) RunServerFunc(ctx context.Context) (err error) {
	errCh := make(chan error)

	go func() {
		errCh <- s.Serve()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.Shutdown()

		return ctx.Err()
	}
}
