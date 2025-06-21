package pubsub

// // PubSubInitHandlers - подключаем обработчики PubSub и SSE
// func (s *Service) PubSubInitHandlers() error {
// 	// command interface handlers
// 	err := s.Routers.CommandProcessor.AddHandlers(
// 		pubsub.NewCommandHandler("GetSmartMessageHandler", s.CommandGetSmartMessageHandler),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("s.Routers.CommandProcessor.AddHandlers: %v", err)
// 	}

// 	// event interface handlers
// 	err = s.Routers.EventProcessor.AddHandlers(
// 		pubsub.NewEventHandler("SmartMessageHandler", s.EventSmartMessageHandler),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("s.Routers.EventProcessor.AddHandlers: %v", err)
// 	}

// 	// sse router handler
// 	s.SSEStream = NewSSEStream()
// 	s.SSEStream.SSEHttpHandlerFunc = s.Routers.SSERouter.AddHandler("events.SmartMessage", s.SSEStream)

// 	return nil
// }

// New инициализирует сущность: сервис приложения
// func New(ctx context.Context, config *config.Config) (*Service, error) {
// 	......

// 	// транспорт PubSub
// 	routers, err := pubsub.NewRouters(ctx, config.PubSub)
// 	if err != nil {
// 		return nil, fmt.Errorf("pubsub.NewRouters:%w", err)
// 	}
//  ......
// }

// // PreServerShutdown - обработчик остановки приложения - вызывается из сервера
// func (s *Service) PreServerShutdown() {
// ......
// 	s.Routers.SSERouter.Close()
// ......
// }

// // OnShutdown - обработчик остановки приложения - вызывается из сервера
// func (s *Service) OnShutdown() {
// 	.......
// 	s.Routers.MessageRouter.Close()
// 	.......
// }
