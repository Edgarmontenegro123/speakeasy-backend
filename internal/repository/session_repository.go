package repository

import (
	"sync"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

type SessionRepository interface {
	Create(session model.Session) model.Session
}

type sessionRepository struct {
	mu       sync.Mutex
	sessions map[string]model.Session
}

func NewSessionRepository() SessionRepository {
	return &sessionRepository{
		sessions: make(map[string]model.Session),
	}
}

func (r *sessionRepository) Create(session model.Session) model.Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return session
}
