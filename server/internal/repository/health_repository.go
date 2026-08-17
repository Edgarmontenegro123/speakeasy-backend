package repository

import "speakeasy-server/internal/model"

type HealthRepository interface {
	GetStatus() model.HealthStatus
}

type healthRepository struct{}

func NewHealthRepository() HealthRepository {
	return &healthRepository{}
}

func (r *healthRepository) GetStatus() model.HealthStatus {
	return model.HealthStatus{Status: "ok"}
}
