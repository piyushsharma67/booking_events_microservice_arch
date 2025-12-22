package service

import "github.com/streadway/amqp"


type Notifier interface {
	SendNotification(to string, subject string, body string) error
}

type MessageBrokerService struct {
	ch  *amqp.Channel
	que string
}

func NewRabbitMQNotifier(conn *amqp.Connection, queue string) (*MessageBrokerService, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &MessageBrokerService{ch: ch, que: queue}, nil
}

func (r *MessageBrokerService) SendNotification(to, subject, body string) error {
    // publish to RabbitMQ
    return nil
}
