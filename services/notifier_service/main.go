package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

// Notification represents the payload sent from auth service
type Notification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

var logger *slog.Logger

func initLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // change via env
	})
	logger = slog.New(handler)
}

func main() {
	initLogger()
	loadEnv()

	// RabbitMQ connection
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		getEnv("RABBITMQ_USER"),
		getEnv("RABBITMQ_PASSWORD"),
		getEnv("RABBITMQ_HOST"),
		getEnv("RABBITMQ_PORT"),
	)
	logger.Info("rabbitmq url",
		"url", rabbitURL,
	)

	var conn *amqp.Connection
	var err error

	for {
		logger.Info("connecting to rabbitmq",
			"host", getEnv("RABBITMQ_HOST"),
		)
		conn, err = amqp.Dial(rabbitURL)
		if err != nil {
			logger.Error("Failed to connect to RabbitMQ, retrying in 5 seconds...", slog.String("error", err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	logger.Info("connected to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open channel", slog.String("error", err.Error()))
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	queueName := getEnv("RABBITMQ_QUEUE")

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
		logger.Error("Failed to declare queue:", slog.String("error", err.Error()))
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
		logger.Error("Failed to register consumer:", slog.String("error", err.Error()))
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)
	logger.Info("Notification service is listening...")

	go func() {
		for d := range msgs {
			msgID := time.Now().UnixNano()
			logger.Info("message received",
				"msg_id", msgID,
				"size", len(d.Body),
				"body", string(d.Body),
			)
			var notif Notification
			if err := json.Unmarshal(d.Body, &notif); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}
			logger.Info("notification received",
				"notification", notif,
			)
			logger.Info("sending email",
				"msg_id", msgID,
				"to", notif.To,
				"subject", notif.Subject,
			)

			start := time.Now()
			duration := time.Since(start)
			if err := sendEmail(notif); err != nil {
				logger.Error("email send failed",
					"msg_id", msgID,
					"to", notif.To,
					"error", err,
					"duration_ms", duration.Milliseconds(),
				)
				continue
			} else {
				logger.Info("email sent successfully",
					"msg_id", msgID,
					"to", notif.To,
					"duration_ms", duration.Milliseconds(),
				)
			}
		}
	}()

	<-forever
}

// sendEmail sends a simple email using SMTP
func sendEmail(n Notification) error {
	smtpHost := getEnv("SMTP_HOST")
	smtpPort := getEnv("SMTP_PORT")
	smtpUser := getEnv("SMTP_USER")
	smtpPass := getEnv("SMTP_PASS")

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		smtpUser, n.To, n.Subject, n.Body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{n.To}, []byte(msg))
}

// getEnv fetches environment variable or returns fallback
func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		return ""
	}
	return val
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		// .env does not exist → do nothing
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Do not override existing env vars
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
