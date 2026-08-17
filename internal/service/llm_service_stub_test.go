package service

import (
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

func TestStubLLMService_GenerateReply(t *testing.T) {
	s := NewStubLLMService()

	reply, err := s.GenerateReply(model.Topic{}, nil, "Hello!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reply != stubTutorReply {
		t.Errorf("expected reply %q, got %q", stubTutorReply, reply)
	}
}
