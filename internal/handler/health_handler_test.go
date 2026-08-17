package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

type stubHealthService struct{}

func (stubHealthService) CheckHealth() model.HealthStatus {
	return model.HealthStatus{Status: "ok"}
}

func TestHealthHandler_GetHealth(t *testing.T) {
	h := NewHealthHandler(stubHealthService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	h.GetHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body model.HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", body.Status)
	}
}
