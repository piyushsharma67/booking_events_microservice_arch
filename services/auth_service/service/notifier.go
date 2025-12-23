package service

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

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

type EmailNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (r *MessageBrokerService) SendNotification(to, subject, body string) error {
	msg := EmailNotification{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		"",    // exchange (default)
		r.que, // routing key = queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // survives broker restart
			Body:         data,
		},
	)
}
