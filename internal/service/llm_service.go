package service

import "github.com/Edgarmontenegro123/speakeasy-backend/internal/model"

// LLMService generates the tutor's reply to a student's message, given the
// active topic and the conversation so far.
type LLMService interface {
	GenerateReply(topic model.Topic, history []model.Message, userMessage string) (string, error)
}
