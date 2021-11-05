package main

import (
	"log"
	"os"

	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/server"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("couldn't load env vars: %v", err)
		}
	}

	DB := &db.PostgresDB{}

	DB.Init()
	s := &server.Server{
		DB:     DB,
		Router: router.NewRouter(),
	}
	s.Start()
}
