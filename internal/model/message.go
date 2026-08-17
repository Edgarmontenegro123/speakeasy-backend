package model

import "time"

type Sender string

const (
	SenderUser Sender = "user"
	SenderAI   Sender = "ai"
)

type Message struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Sender    Sender    `json:"sender"`
	Content   string    `json:"content"`
	AudioURL  string    `json:"audio_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
