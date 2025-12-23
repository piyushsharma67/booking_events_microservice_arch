package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/streadway/amqp"
)

// Notification represents the payload sent from auth service
type Notification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	// RabbitMQ connection
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		getEnv("RABBITMQ_USER", "guest"),
		getEnv("RABBITMQ_PASSWORD", "guest"),
		getEnv("RABBITMQ_HOST", "rabbitmq"),
		getEnv("RABBITMQ_PORT", "5672"),
	)

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	queueName := getEnv("RABBITMQ_QUEUE", "notifications")

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)
	log.Println("Notification service is listening...")

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			var notif Notification
			if err := json.Unmarshal(d.Body, &notif); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			if err := sendEmail(notif); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email sent to %s", notif.To)
			}
		}
	}()

	<-forever
}

// sendEmail sends a simple email using SMTP
func sendEmail(n Notification) error {
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "your-email@gmail.com")
	smtpPass := getEnv("SMTP_PASS", "your-password")

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		smtpUser, n.To, n.Subject, n.Body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{n.To}, []byte(msg))
}

// getEnv fetches environment variable or returns fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
