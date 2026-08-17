package handler

import "net/http"

// placeholderAudio is a stand-in payload served for generated tutor audio
// until real text-to-speech synthesis is integrated.
var placeholderAudio = []byte("stub audio content")

type AudioHandler struct{}

func NewAudioHandler() *AudioHandler {
	return &AudioHandler{}
}

func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(placeholderAudio)
}
