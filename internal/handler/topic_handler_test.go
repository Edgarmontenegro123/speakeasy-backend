package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

type stubTopicService struct {
	topics []model.Topic
}

func (s stubTopicService) ListTopics() []model.Topic {
	return s.topics
}

func TestTopicHandler_ListTopics(t *testing.T) {
	want := []model.Topic{
		{ID: "topic-ordering-coffee", Title: "Ordering Coffee", Description: "Practise ordering a coffee at a café.", Level: model.LevelA1},
	}
	h := NewTopicHandler(stubTopicService{topics: want})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
	rec := httptest.NewRecorder()

	h.ListTopics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []model.Topic
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(got) != len(want) || got[0].ID != want[0].ID {
		t.Errorf("expected topics %+v, got %+v", want, got)
	}
}
