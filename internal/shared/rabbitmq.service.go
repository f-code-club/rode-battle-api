package shared

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpService struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewRabbitService(url string, taskQueue string) (*AmqpService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		taskQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		return nil, err
	}

	return &AmqpService{
		Channel: ch,
		Queue:   q,
	}, nil
}
