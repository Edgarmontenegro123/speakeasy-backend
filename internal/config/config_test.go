package config

import "testing"

func TestLoad_DefaultsGeminiModel(t *testing.T) {
	t.Setenv("GEMINI_MODEL", "")

	cfg := Load()

	if cfg.GeminiModel != defaultGeminiModel {
		t.Errorf("expected default Gemini model %q, got %q", defaultGeminiModel, cfg.GeminiModel)
	}
}

func TestLoad_UsesConfiguredGeminiModel(t *testing.T) {
	t.Setenv("GEMINI_MODEL", "gemini-2.5-pro")

	cfg := Load()

	if cfg.GeminiModel != "gemini-2.5-pro" {
		t.Errorf("expected Gemini model %q, got %q", "gemini-2.5-pro", cfg.GeminiModel)
	}
}
