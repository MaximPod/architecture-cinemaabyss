package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

// NewEventKafka - инициализирует канал на базе Kafka для событийного интерфейса
// возвращает Publisher и функцию создания Subscribers
func NewEventKafka(logger watermill.LoggerAdapter, cfg KafkaConfig) (message.Publisher, func(topicName string) (message.Subscriber, error)) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{cfg.Brokers},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	subscriberFunc := func(topicName string) (message.Subscriber, error) {

		kafkaSubscriber, err := kafka.NewSubscriber(
			kafka.SubscriberConfig{
				Brokers:       []string{cfg.Brokers},
				Unmarshaler:   kafka.DefaultMarshaler{},
				ConsumerGroup: generateEventsTopic(topicName), // every handler will use a separate consumer group
			},
			logger,
		)
		if err != nil {
			panic(err)
		}

		return kafkaSubscriber, err
	}

	return kafkaPublisher, subscriberFunc
}
