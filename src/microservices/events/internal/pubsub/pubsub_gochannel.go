package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// NewEventGoChannel - инициализирует канал на базе gochannel для событийного интерфейса
// возвращает Publisher и функцию создания Subscribers
func NewEventGoChannel(logger watermill.LoggerAdapter) (message.Publisher, func(topicName string) (message.Subscriber, error)) {
	goChannel := gochannel.NewGoChannel(
		gochannel.Config{},
		logger,
	)

	subscriberFunc := func(topicName string) (message.Subscriber, error) {
		return goChannel, nil
	}

	return goChannel, subscriberFunc
}
