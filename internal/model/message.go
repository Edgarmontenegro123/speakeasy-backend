package model

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	AudioURL  string    `json:"audio_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
