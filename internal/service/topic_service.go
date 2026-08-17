package service

import (
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
)

type TopicService interface {
	ListTopics() []model.Topic
}

type topicService struct {
	repo repository.TopicRepository
}

func NewTopicService(repo repository.TopicRepository) TopicService {
	return &topicService{repo: repo}
}

func (s *topicService) ListTopics() []model.Topic {
	return s.repo.ListTopics()
}
