package messaging

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type UserRegisteredEventHandler func(event UserRegisteredEvent) error

type NatsSubscriber struct {
	natsConnection *nats.Conn
}

func NewNatsSubscriber(natsURL string) (*NatsSubscriber, error) {
	natsConnection, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &NatsSubscriber{natsConnection: natsConnection}, nil
}

func (s *NatsSubscriber) SubscribeToUserRegisteredEvents(handler UserRegisteredEventHandler) error {
	_, err := s.natsConnection.Subscribe("user.registered", func(message *nats.Msg) {
		var event UserRegisteredEvent
		if err := json.Unmarshal(message.Data, &event); err != nil {
			log.Printf("failed to unmarshal user registered event: %v", err)
			return
		}

		if err := handler(event); err != nil {
			log.Printf("failed to handle user registered event: %v", err)
			return
		}

		log.Printf("processed user registered event for user: %s", event.UserID)
	})
	return err
}

func (s *NatsSubscriber) Close() {
	s.natsConnection.Close()
}
