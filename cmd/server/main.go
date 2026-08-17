package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := router.New()

	log.Printf("Speakeasy API listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
