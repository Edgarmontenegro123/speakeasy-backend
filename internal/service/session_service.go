package service

import (
	"errors"
	"time"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

var (
	ErrTopicNotFound   = errors.New("topic not found")
	ErrSessionNotFound = errors.New("session not found")
)

const tutorReply = "Great start! Can you tell me more about that?"

type SessionService interface {
	CreateSession(topicID string) (model.Session, error)
	PostMessage(sessionID, content string) (model.Message, error)
}

type sessionService struct {
	sessions repository.SessionRepository
	topics   repository.TopicRepository
	tts      TTSService
}

func NewSessionService(sessions repository.SessionRepository, topics repository.TopicRepository, tts TTSService) SessionService {
	return &sessionService{sessions: sessions, topics: topics, tts: tts}
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

	userMessage := model.Message{
		ID:        newID("message"),
		SessionID: sessionID,
		Sender:    model.SenderUser,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	tutorMessage := model.Message{
		ID:        newID("message"),
		SessionID: sessionID,
		Sender:    model.SenderAI,
		Content:   tutorReply,
		CreatedAt: time.Now().UTC(),
	}
	tutorMessage.AudioURL = s.tts.GenerateAudioURL(tutorMessage.ID)

	session.Messages = append(session.Messages, userMessage, tutorMessage)
	s.sessions.Update(session)

	return tutorMessage, nil
}

func (s *sessionService) topicExists(topicID string) bool {
	for _, topic := range s.topics.ListTopics() {
		if topic.ID == topicID {
			return true
		}
	}
	return false
}
