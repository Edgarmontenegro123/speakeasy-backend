package main

import (
	"log"
	"net/http"

	"speakeasy-server/internal/config"
	"speakeasy-server/internal/router"
)

func main() {
	cfg := config.Load()
	mux := router.New()

	log.Printf("Speakeasy API listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
