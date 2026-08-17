package service

import "fmt"

// TTSService turns tutor message text into audio. This stub implementation
// does not perform real speech synthesis; it returns the URL where the
// generated audio for a message would be served from.
type TTSService interface {
	GenerateAudioURL(messageID string) string
}

type ttsService struct{}

func NewTTSService() TTSService {
	return &ttsService{}
}

func (s *ttsService) GenerateAudioURL(messageID string) string {
	return fmt.Sprintf("/api/v1/audio/%s.mp3", messageID)
}
