package config

import "os"

type Config struct {
	Port         string
	GeminiAPIKey string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port:         port,
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}
