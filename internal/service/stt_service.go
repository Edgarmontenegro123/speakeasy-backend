package service

// STTService turns spoken audio into text. This stub implementation does not
// perform real speech recognition; it returns a simulated transcription
// regardless of the audio content provided.
type STTService interface {
	Transcribe(audio []byte, mimeType string) (string, error)
}

type sttService struct{}

func NewSTTService() STTService {
	return &sttService{}
}

func (s *sttService) Transcribe(audio []byte, mimeType string) (string, error) {
	return "Hello, I want to practice my English today", nil
}
