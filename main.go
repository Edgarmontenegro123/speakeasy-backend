package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/config"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using existing environment variables")
	}

	cfg := config.Load()
	mux := router.New(cfg)

	log.Printf("Speakeasy API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
