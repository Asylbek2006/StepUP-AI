package messaging

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type NatsPublisher struct {
	natsConnection *nats.Conn
}

func NewNatsPublisher(natsURL string) (*NatsPublisher, error) {
	natsConnection, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &NatsPublisher{natsConnection: natsConnection}, nil
}

func (p *NatsPublisher) PublishUserRegisteredEvent(event UserRegisteredEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.natsConnection.Publish("user.registered", eventData)
}

func (p *NatsPublisher) Close() {
	p.natsConnection.Close()
}
