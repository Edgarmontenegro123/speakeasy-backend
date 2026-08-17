package model

import "time"

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
)

type Session struct {
	ID        string        `json:"id"`
	TopicID   string        `json:"topic_id"`
	Status    SessionStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Messages  []Message     `json:"messages"`
}
