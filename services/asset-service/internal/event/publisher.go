package event

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

type EventPublisher interface {
	Publish(subject string, v any) error
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) Publish(subject string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.nc.Publish(subject, data)
}
