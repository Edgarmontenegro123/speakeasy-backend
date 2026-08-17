package service

import (
	"errors"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

type mockLLMService struct {
	reply      string
	err        error
	gotTopic   model.Topic
	gotHistory []model.Message
	gotMessage string
}

func (m *mockLLMService) GenerateReply(topic model.Topic, history []model.Message, userMessage string) (string, error) {
	m.gotTopic = topic
	m.gotHistory = history
	m.gotMessage = userMessage

	if m.err != nil {
		return "", m.err
	}
	return m.reply, nil
}

func TestSessionService_CreateSession(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService(), NewStubLLMService())

	topicID := topics.ListTopics()[0].ID

	session, err := s.CreateSession(topicID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.TopicID != topicID {
		t.Errorf("expected topic ID %q, got %q", topicID, session.TopicID)
	}

	if session.ID == "" {
		t.Error("expected a generated session ID, got empty string")
	}

	if session.Status != "active" {
		t.Errorf("expected status %q, got %q", "active", session.Status)
	}

	if session.Messages == nil {
		t.Error("expected Messages to be initialised, got nil")
	}
}

func TestSessionService_CreateSession_UnknownTopic(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService(), NewStubLLMService())

	_, err := s.CreateSession("unknown-topic")
	if !errors.Is(err, ErrTopicNotFound) {
		t.Fatalf("expected ErrTopicNotFound, got %v", err)
	}
}

func TestSessionService_PostMessage(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService(), NewStubLLMService())

	session, err := s.CreateSession(topics.ListTopics()[0].ID)
	if err != nil {
		t.Fatalf("unexpected error creating session: %v", err)
	}

	reply, err := s.PostMessage(session.ID, "Hello, I want to practise English.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply.Sender != model.SenderAI {
		t.Errorf("expected sender %q, got %q", model.SenderAI, reply.Sender)
	}

	if reply.Content == "" {
		t.Error("expected a non-empty tutor reply")
	}

	if reply.AudioURL == "" {
		t.Error("expected a non-empty audio URL for the tutor reply")
	}

	updated, ok := sessions.Get(session.ID)
	if !ok {
		t.Fatal("expected session to exist after posting a message")
	}

	if len(updated.Messages) != 2 {
		t.Fatalf("expected 2 messages stored, got %d", len(updated.Messages))
	}

	if updated.Messages[0].Sender != model.SenderUser || updated.Messages[0].Content != "Hello, I want to practise English." {
		t.Errorf("unexpected first stored message: %+v", updated.Messages[0])
	}

	if updated.Messages[1].Sender != model.SenderAI {
		t.Errorf("unexpected second stored message sender: %q", updated.Messages[1].Sender)
	}

	if updated.Messages[1].AudioURL == "" {
		t.Error("expected the stored tutor message to have a non-empty audio URL")
	}
}

func TestSessionService_PostMessage_SessionNotFound(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService(), NewStubLLMService())

	_, err := s.PostMessage("unknown-session", "Hello")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestSessionService_PostMessage_UsesLLMService(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	llm := &mockLLMService{reply: "Nice! Tell me about your favourite coffee shop."}
	s := NewSessionService(sessions, topics, NewTTSService(), llm)

	topic := topics.ListTopics()[0]
	session, err := s.CreateSession(topic.ID)
	if err != nil {
		t.Fatalf("unexpected error creating session: %v", err)
	}

	reply, err := s.PostMessage(session.ID, "I had a coffee this morning.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply.Content != llm.reply {
		t.Errorf("expected reply content %q, got %q", llm.reply, reply.Content)
	}

	if llm.gotTopic.ID != topic.ID {
		t.Errorf("expected the LLM service to receive topic %q, got %q", topic.ID, llm.gotTopic.ID)
	}

	if llm.gotMessage != "I had a coffee this morning." {
		t.Errorf("expected the LLM service to receive the user's message, got %q", llm.gotMessage)
	}
}

func TestSessionService_PostMessage_LLMError(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	llm := &mockLLMService{err: errors.New("gemini unavailable")}
	s := NewSessionService(sessions, topics, NewTTSService(), llm)

	session, err := s.CreateSession(topics.ListTopics()[0].ID)
	if err != nil {
		t.Fatalf("unexpected error creating session: %v", err)
	}

	if _, err := s.PostMessage(session.ID, "Hello"); err == nil {
		t.Fatal("expected an error when the LLM service fails")
	}

	updated, _ := sessions.Get(session.ID)
	if len(updated.Messages) != 0 {
		t.Errorf("expected no messages stored when the LLM call fails, got %d", len(updated.Messages))
	}
}
