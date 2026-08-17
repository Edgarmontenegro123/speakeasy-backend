package service

import (
	"errors"
	"time"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

var ErrTopicNotFound = errors.New("topic not found")

type SessionService interface {
	CreateSession(topicID string) (model.Session, error)
}

type sessionService struct {
	sessions repository.SessionRepository
	topics   repository.TopicRepository
}

func NewSessionService(sessions repository.SessionRepository, topics repository.TopicRepository) SessionService {
	return &sessionService{sessions: sessions, topics: topics}
}

func (s *sessionService) CreateSession(topicID string) (model.Session, error) {
	if !s.topicExists(topicID) {
		return model.Session{}, ErrTopicNotFound
	}

	session := model.Session{
		ID:        newID("session"),
		TopicID:   topicID,
		Status:    model.SessionStatusActive,
		CreatedAt: time.Now().UTC(),
		Messages:  []model.Message{},
	}

	return s.sessions.Create(session), nil
}

func (s *sessionService) topicExists(topicID string) bool {
	for _, topic := range s.topics.ListTopics() {
		if topic.ID == topicID {
			return true
		}
	}
	return false
}
