package rabbitmq

import (
	"fmt"
	"time"

	"github.com/benu-cloud/benu-livestreaming-gst/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	conn           *amqp.Connection
	publishTimeout time.Duration
}

func NewConnection(config *config.MessageBrokerSettings) (*RabbitMQConnection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VHost)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConnection{conn: conn, publishTimeout: config.PublishTimeout}, nil
}

func (r *RabbitMQConnection) CloseConnection() error {
	return r.conn.Close()
}
