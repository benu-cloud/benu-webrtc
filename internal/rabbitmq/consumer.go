package rabbitmq

import (
	pkgerrors "github.com/benu-cloud/benu-errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQConnection) NewConsumer(queuename string) (<-chan amqp.Delivery, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, pkgerrors.NewConsumeError(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queuename,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, pkgerrors.NewConsumeError(err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, pkgerrors.NewConsumeError(err)
	}
	return msgs, nil
}
