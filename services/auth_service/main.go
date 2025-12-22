package main

import (
	"log"
	"net/http"

	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/repository"
	"github.com/piyushsharma67/movie_booking/services/auth_service/routes"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
)

func main() {

	// 1️⃣ Initialize low-level DB (needs Close)
	pgxpool, queries := databases.InitPostgres()
	defer pgxpool.Close()

	// 2️⃣ Wrap with interface
	db := databases.NewPostgresDB(queries)
	repository := repository.NewUserRepository(db)

	srv:=service.NewAuthService(repository)
	r:=routes.InitRoutes(srv)

	log.Println("Server running on :8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}
