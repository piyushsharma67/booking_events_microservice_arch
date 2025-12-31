package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piyushsharma67/events_booking/services/events_service/database"
	"github.com/piyushsharma67/events_booking/services/events_service/logger"
	"github.com/piyushsharma67/events_booking/services/events_service/repository"
	"github.com/piyushsharma67/events_booking/services/events_service/routes"
	"github.com/piyushsharma67/events_booking/services/events_service/service"
)

func main() {

	logger := logger.NewSlogFileLogger("events", "development", "./logs/events/events.log", slog.LevelInfo)
	logger.Info(fmt.Sprintf("events Server running on port :%s", os.Getenv("SERVER_PORT")))

	/* genrating the db connection */
	// 1️⃣ Initialize low-level DB (needs Close)
	// pgxpool, queries := database.InitPostgres()
	// defer pgxpool.Close()

	mongodbClient, close := database.ConnectMongo()
	defer close()

	db := database.NewMongoDb(mongodbClient)

	// db:=database.NewSqldb(queries)

	repository := repository.NewRepos(db)

	srv := service.GetEventService(*repository)

	r := routes.InitRoutes(srv, logger)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8003"
	}

	httpServer := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	// ---------- START SERVER ----------
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	logger.Info("Waiting for the termination signal..")
	<-quit
	logger.Info("terminate signal recieved")

	// ---------- GRACEFUL SHUTDOWN ----------
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("server forced shutdown", "error", err.Error())
	} else {
		logger.Info("server shutdown gracefully")
	}

}
