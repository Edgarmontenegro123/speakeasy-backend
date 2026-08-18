package service

import (
	"strings"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

func TestBuildTutorPrompt(t *testing.T) {
	topic := model.Topic{
		ID:          "topic-ordering-coffee",
		Title:       "Ordering Coffee",
		Description: "Practise ordering a coffee at a café.",
		Level:       model.LevelA1,
	}
	history := []model.Message{
		{Role: model.RoleUser, Content: "Hi!"},
		{Role: model.RoleAssistant, Content: "Hello! Ready to practise?"},
	}

	prompt := buildTutorPrompt(topic, history, "I would like a coffee.")

	for _, want := range []string{
		string(topic.Level),
		topic.Title,
		topic.Description,
		"Student: Hi!",
		"Tutor: Hello! Ready to practise?",
		"Student: I would like a coffee.",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("expected prompt to contain %q, got:\n%s", want, prompt)
		}
	}
}

func TestLevelGuidance(t *testing.T) {
	cases := map[model.Level]string{
		model.LevelA1: "simple",
		model.LevelA2: "simple",
		model.LevelB1: "everyday",
		model.LevelB2: "everyday",
		model.LevelC1: "idiomatic",
		model.LevelC2: "idiomatic",
	}

	for level, want := range cases {
		guidance := levelGuidance(level)
		if !strings.Contains(guidance, want) {
			t.Errorf("expected guidance for level %q to contain %q, got %q", level, want, guidance)
		}
	}
}
