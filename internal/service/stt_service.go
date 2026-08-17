package service

// STTService turns spoken audio into text. This stub implementation does not
// perform real speech recognition; it returns a simulated transcription
// regardless of the audio content provided.
type STTService interface {
	Transcribe(audio []byte) string
}

type sttService struct{}

func NewSTTService() STTService {
	return &sttService{}
}

func (s *sttService) Transcribe(audio []byte) string {
	return "Hello, I want to practice my English today"
}
