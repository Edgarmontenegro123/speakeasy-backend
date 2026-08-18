package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

var (
	ErrTopicNotFound   = errors.New("topic not found")
	ErrSessionNotFound = errors.New("session not found")
)

type SessionService interface {
	CreateSession(topicID string) (model.Session, error)
	PostMessage(sessionID, content string) (model.Message, error)
}

type sessionService struct {
	sessions repository.SessionRepository
	topics   repository.TopicRepository
	tts      TTSService
	llm      LLMService
}

func NewSessionService(sessions repository.SessionRepository, topics repository.TopicRepository, tts TTSService, llm LLMService) SessionService {
	return &sessionService{sessions: sessions, topics: topics, tts: tts, llm: llm}
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

func (s *sessionService) PostMessage(sessionID, content string) (model.Message, error) {
	session, ok := s.sessions.Get(sessionID)
	if !ok {
		return model.Message{}, ErrSessionNotFound
	}

	topic, ok := s.topics.GetTopic(session.TopicID)
	if !ok {
		return model.Message{}, ErrTopicNotFound
	}

	reply, err := s.llm.GenerateReply(topic, session.Messages, content)
	if err != nil {
		return model.Message{}, fmt.Errorf("generating tutor reply: %w", err)
	}

	userMessage := model.Message{
		ID:        newID("message"),
		SessionID: sessionID,
		Role:      model.RoleUser,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	tutorMessage := model.Message{
		ID:        newID("message"),
		SessionID: sessionID,
		Role:      model.RoleAssistant,
		Content:   reply,
		CreatedAt: time.Now().UTC(),
	}
	tutorMessage.AudioURL = s.tts.GenerateAudioURL(tutorMessage.ID)

	session.Messages = append(session.Messages, userMessage, tutorMessage)
	s.sessions.Update(session)

	return tutorMessage, nil
}

func (s *sessionService) topicExists(topicID string) bool {
	_, ok := s.topics.GetTopic(topicID)
	return ok
}
