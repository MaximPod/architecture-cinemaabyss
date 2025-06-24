package pubsub

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// PubSubConfig - конфигурация pubsub
type PubSubConfig struct {
	TransportType string      `yaml:"transport_type"` // Types: gochannel, kafka
	Kafka         KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers string `yaml:"brokers"` // список брокеров
	Topic   string `yaml:"topic"`   // имя топика
}

// EventPublisher - интерфейс публикации сообщений
type EventPublisher interface {
	Publish(eventType string, events ...*message.Message) error
}

// Routers - транспортный модуль pubsub модели взаимодействия
type Routers struct {
	MessageRouter  *message.Router      // маршрутизатор сообщений
	EventBus       *cqrs.EventBus       // шина событий
	EventProcessor *cqrs.EventProcessor // обработчик событий из шины
}

// NewRouters - конуструктор транспортной системы
func NewRouters(ctx context.Context, cfg PubSubConfig) (*Routers, error) {
	logger := watermill.NewSlogLogger(nil)

	// message router
	messageRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("messageRouter.NewRouter: %v", err)
	}

	messageRouter.AddMiddleware(middleware.Recoverer)

	var eventsPublisher message.Publisher
	var eventsSubscriberFn func(topicName string) (message.Subscriber, error)

	switch cfg.TransportType {
	case "gochannel", "":
		eventsPublisher, eventsSubscriberFn = NewEventGoChannel(logger)

	case "kafka":
		eventsPublisher, eventsSubscriberFn = NewEventKafka(logger, cfg.Kafka)
	// case "oneOfTheListBelow":
	// так же поддерживаются транспорты типа:
	// AMQP 1.0
	// Apache Pulsar
	// Apache RocketMQ
	// CockroachDB
	// Ensign
	// GoogleCloud Pub/Sub HTTP Push
	// MongoDB
	// MQTT
	// NSQ
	// Redis Zset
	// SQLite
	// Postgres
	// https://watermill.io/docs/awesome/
	default:
		return nil, fmt.Errorf("unknown transport type %s", cfg.TransportType)
	}

	eventBus, eventProcessor, err := newEventsInterface(messageRouter, eventsPublisher, eventsSubscriberFn, logger)
	if err != nil {
		return nil, fmt.Errorf("newCommandsInterface: %v", err)
	}

	routers := Routers{
		MessageRouter:  messageRouter,
		EventBus:       eventBus,
		EventProcessor: eventProcessor,
	}

	return &routers, nil
}
