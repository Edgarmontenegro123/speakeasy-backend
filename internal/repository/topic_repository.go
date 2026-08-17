package repository

import "github.com/Edgarmontenegro123/speakeasy-backend/internal/model"

type TopicRepository interface {
	ListTopics() []model.Topic
	GetTopic(id string) (model.Topic, bool)
}

type topicRepository struct {
	topics []model.Topic
}

func NewTopicRepository() TopicRepository {
	return &topicRepository{
		topics: []model.Topic{
			{ID: "topic-ordering-coffee", Title: "Ordering Coffee", Description: "Practise ordering a coffee at a café.", Level: model.LevelA1},
			{ID: "topic-introducing-yourself", Title: "Introducing Yourself", Description: "Practise a first conversation with a new colleague.", Level: model.LevelA2},
			{ID: "topic-travel-planning", Title: "Travel Planning", Description: "Practise discussing travel plans and itineraries.", Level: model.LevelB1},
			{ID: "topic-job-interview", Title: "Job Interview", Description: "Practise answering common job interview questions.", Level: model.LevelB2},
			{ID: "topic-debating-ideas", Title: "Debating Ideas", Description: "Practise defending an opinion in a debate.", Level: model.LevelC1},
		},
	}
}

func (r *topicRepository) ListTopics() []model.Topic {
	return r.topics
}

func (r *topicRepository) GetTopic(id string) (model.Topic, bool) {
	for _, topic := range r.topics {
		if topic.ID == id {
			return topic, true
		}
	}
	return model.Topic{}, false
}
