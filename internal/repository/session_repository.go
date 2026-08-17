package repository

import (
	"sync"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

type SessionRepository interface {
	Create(session model.Session) model.Session
	Get(id string) (model.Session, bool)
	Update(session model.Session) (model.Session, bool)
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

func (r *sessionRepository) Get(id string) (model.Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[id]
	return session, ok
}

func (r *sessionRepository) Update(session model.Session) (model.Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.ID]; !ok {
		return model.Session{}, false
	}

	r.sessions[session.ID] = session
	return session, true
}
