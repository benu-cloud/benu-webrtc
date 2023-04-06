package rabbitmq

import (
	"context"

	"github.com/benu-cloud/benu-webrtc/internal/message"
	"github.com/benu-cloud/benu-webrtc/pkg/pkgerrors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQConnection) Publish(queuename string, payload message.GenericPayload) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return pkgerrors.NewPublishError(err)
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
		return pkgerrors.NewPublishError(err)
	}

	body, err := message.Marshal(payload)
	if err != nil {
		return pkgerrors.NewPublishError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.publishTimeout)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // default exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return pkgerrors.NewPublishError(err)
	}
	return nil
}
