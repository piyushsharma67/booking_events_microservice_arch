package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/piyushsharma67/events_booking/services/auth_service/databases"
	"github.com/piyushsharma67/events_booking/services/auth_service/logger"
	"github.com/piyushsharma67/events_booking/services/auth_service/repository"
	"github.com/piyushsharma67/events_booking/services/auth_service/routes"
	"github.com/piyushsharma67/events_booking/services/auth_service/service"
	"github.com/streadway/amqp"
)

func main() {
	logger := logger.NewSlogFileLogger("auth_service", "development", "./logs/auth_service/auth.log", slog.LevelInfo)

	dbClient,close:=databases.ConnectMongo()
	defer close()
	db:=databases.NewMongoDb(dbClient)
	logger.Info("db connected")
	
	repository := repository.NewUserRepository(db)
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	var conn *amqp.Connection
	var err error

	rabbitmqDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	// Retry logic for RabbitMQ
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(rabbitmqDSN)
		if err == nil {
			break
		}
		slog.Warn("RabbitMQ not ready, retrying...", "attempt", i+1, "error", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ failed to connect to rabbitmq after retries: %v", err)
	}
	defer conn.Close()
	
	notifier, err := service.NewRabbitMQNotifier(conn, "notifications")
	if err != nil {
		logger.Error("Error occured for auth service", "error", err.Error())
		log.Fatal(err)

	}
	srv := service.NewAuthService(repository, notifier, logger)
	r := routes.InitRoutes(srv, logger)

	log.Println("Server running on :8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}
