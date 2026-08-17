package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

const geminiModel = "gemini-1.5-flash"

// geminiLLMService generates tutor replies using the Gemini API.
type geminiLLMService struct {
	client *genai.Client
}

func NewGeminiLLMService(ctx context.Context, apiKey string) (LLMService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &geminiLLMService{client: client}, nil
}

func (s *geminiLLMService) GenerateReply(topic model.Topic, history []model.Message, userMessage string) (string, error) {
	prompt := buildTutorPrompt(topic, history, userMessage)

	resp, err := s.client.Models.GenerateContent(context.Background(), geminiModel, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generating tutor reply: %w", err)
	}

	reply := strings.TrimSpace(resp.Text())
	if reply == "" {
		return "", fmt.Errorf("gemini returned an empty reply")
	}

	return reply, nil
}

// buildTutorPrompt assembles a prompt that grounds the model in the active
// topic, adapts its language to the topic's CEFR level, and includes the
// conversation so far so the tutor can respond in context.
func buildTutorPrompt(topic model.Topic, history []model.Message, userMessage string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "You are a friendly English tutor helping a %s-level student practise the topic %q (%s). ",
		topic.Level, topic.Title, topic.Description)
	b.WriteString(levelGuidance(topic.Level))
	b.WriteString(" Keep your reply short (1-3 sentences), gently correct any mistakes, and end with a question that keeps the student talking.\n\n")

	for _, msg := range history {
		switch msg.Sender {
		case model.SenderUser:
			b.WriteString("Student: " + msg.Content + "\n")
		case model.SenderAI:
			b.WriteString("Tutor: " + msg.Content + "\n")
		}
	}

	b.WriteString("Student: " + userMessage + "\nTutor:")

	return b.String()
}

func levelGuidance(level model.Level) string {
	switch level {
	case model.LevelA1, model.LevelA2:
		return "Use very simple vocabulary, short sentences, and the present tense where possible."
	case model.LevelB1, model.LevelB2:
		return "Use everyday vocabulary and moderately complex sentences."
	default:
		return "Use natural, idiomatic English and nuanced vocabulary."
	}
}
