package service

import (
	"speakeasy-server/internal/model"
	"speakeasy-server/internal/repository"
)

type HealthService interface {
	CheckHealth() model.HealthStatus
}

type healthService struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) HealthService {
	return &healthService{repo: repo}
}

func (s *healthService) CheckHealth() model.HealthStatus {
	return s.repo.GetStatus()
}
