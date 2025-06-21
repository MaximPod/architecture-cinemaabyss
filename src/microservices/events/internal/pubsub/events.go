package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

// generateEventsTopic -  генерирует имя топика для каналов с событиями
func generateEventsTopic(eventName string) string {
	return "events." + eventName
}

// newEventsInterface - создает и возвращает событийные шину и процессор
func newEventsInterface(router *message.Router,
	eventsPublisher message.Publisher,
	eventsSubscriberFn func(topicName string) (message.Subscriber, error),
	logger watermill.LoggerAdapter) (*cqrs.EventBus, *cqrs.EventProcessor, error) {

	eventBus, err := cqrs.NewEventBusWithConfig(eventsPublisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return generateEventsTopic(params.EventName), nil
		},

		OnPublish: func(params cqrs.OnEventSendParams) error {
			logger.Debug("Event Bus: event recieved", watermill.LogFields{
				"event_name": params.EventName,
			})

			params.Message.Metadata.Set("published_at", time.Now().String())

			return nil
		},

		Marshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
		Logger: logger,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cqrs.NewEventBusWithConfig: %v", err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			AckOnUnknownEvent: true,
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return generateEventsTopic(params.EventName), nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return eventsSubscriberFn(params.EventName)
			},

			OnHandle: func(params cqrs.EventProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Event)

				logger.Debug("Event Processor: event handled", watermill.LogFields{
					"event_name": params.EventName,
					"duration":   time.Since(start),
					"err":        err,
				})

				return err
			},

			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
			Logger: logger,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cqrs.NewEventProcessorWithConfig: %v", err)
	}

	return eventBus, eventProcessor, nil
}

// NewEventHandler - обертка для регистрации обработчиков eventProcessor
func NewEventHandler[T any](handlerName string, handleFunc func(ctx context.Context, cmd *T) error) cqrs.EventHandler {
	return cqrs.NewEventHandler(handlerName, handleFunc)
}
