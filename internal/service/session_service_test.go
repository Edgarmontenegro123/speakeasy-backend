package service

import (
	"errors"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

func TestSessionService_CreateSession(t *testing.T) {
	topics := repository.NewTopicRepository()
	sessions := repository.NewSessionRepository()
	s := NewSessionService(sessions, topics)

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
	s := NewSessionService(sessions, topics)

	_, err := s.CreateSession("unknown-topic")
	if !errors.Is(err, ErrTopicNotFound) {
		t.Fatalf("expected ErrTopicNotFound, got %v", err)
	}
}
