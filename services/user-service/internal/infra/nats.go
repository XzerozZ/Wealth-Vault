package infra

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func NewNATSConnection(host, port string) (*nats.Conn, error) {
	url := fmt.Sprintf("nats://%s:%s", host, port)

	opts := []nats.Option{
		nats.Name("User-Service"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),

		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("⚠️ NATS Disconnected: %v", err)
			}
		}),

		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("♻️ NATS Reconnected to: %s", nc.ConnectedUrl())
		}),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	return nc, nil
}
