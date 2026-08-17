package service

import (
	"errors"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

func TestSessionService_CreateSession(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService())

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
	s := NewSessionService(sessions, topics, NewTTSService())

	_, err := s.CreateSession("unknown-topic")
	if !errors.Is(err, ErrTopicNotFound) {
		t.Fatalf("expected ErrTopicNotFound, got %v", err)
	}
}

func TestSessionService_PostMessage(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics, NewTTSService())

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
	s := NewSessionService(sessions, topics, NewTTSService())

	_, err := s.PostMessage("unknown-session", "Hello")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound, got %v", err)
	}
}
