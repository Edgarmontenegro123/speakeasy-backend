package service

import "github.com/Edgarmontenegro123/speakeasy-backend/internal/model"

const stubTutorReply = "Great start! Can you tell me more about that?"

// stubLLMService is a canned LLMService used when no Gemini API key is
// configured, and by unit tests that need a deterministic tutor reply.
type stubLLMService struct{}

func NewStubLLMService() LLMService {
	return &stubLLMService{}
}

func (s *stubLLMService) GenerateReply(topic model.Topic, history []model.Message, userMessage string) (string, error) {
	return stubTutorReply, nil
}
