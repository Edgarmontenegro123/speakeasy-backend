package config

import "os"

const defaultGeminiModel = "gemini-flash-latest"

type Config struct {
	Port         string
	GeminiAPIKey string
	GeminiModel  string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		geminiModel = defaultGeminiModel
	}

	return Config{
		Port:         port,
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiModel:  geminiModel,
	}
}
