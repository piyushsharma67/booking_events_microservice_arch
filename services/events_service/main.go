package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piyushsharma67/events_booking/services/events_service/database"
	"github.com/piyushsharma67/events_booking/services/events_service/infra"
	"github.com/piyushsharma67/events_booking/services/events_service/logger"
	"github.com/piyushsharma67/events_booking/services/events_service/que"
	"github.com/piyushsharma67/events_booking/services/events_service/repository"
	"github.com/piyushsharma67/events_booking/services/events_service/routes"
	"github.com/piyushsharma67/events_booking/services/events_service/service"
	"github.com/piyushsharma67/events_booking/services/events_service/utils"
	"github.com/streadway/amqp"
)

func main() {

	logger := logger.NewSlogFileLogger(
		"events",
		"development",
		"./logs/events/events.log",
		slog.LevelInfo,
	)

	logger.Info("events Server starting", "port", os.Getenv("SERVER_PORT"))

	// ---------- MongoDB ----------
	mongodbClient, close := database.ConnectMongo()
	defer close()

	db := database.NewMongoDb(mongodbClient)
	repo := repository.NewRepos(db)

	// ---------- RabbitMQ ----------
	rabbitURL := utils.BuildRabbitURL()
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var conn *amqp.Connection
	var err error

	for i := 1; i <= 20; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			log.Println("RabbitMQ connected")
			break
		}
		log.Printf("RabbitMQ not ready (%d/20): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}

	// ---------- Publisher ----------
	publisher, err := que.NewEventsPublisher(
		conn,
		os.Getenv("RABBITMQ_QUEUE_SEATS"),
	)
	if err != nil {
		log.Fatal(err)
	}

	redisClient:= infra.NewRedisClient()
	// ---------- Service ----------
	srv := service.GetEventService(*repo, publisher,redisClient)

	// ---------- HTTP ----------
	r := routes.InitRoutes(srv, logger)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8003"
	}

	httpServer := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
		}
	}()

	// ---------- Graceful Shutdown ----------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	conn.Close()

	logger.Info("server stopped gracefully")
}
