package config

import "os"

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topics  struct {
		TaskCreated   string
		TaskUpdated   string
		TaskDeleted   string
		TaskOverdue   string
		Notifications string
	}
}

func LoadKafkaConfig() *KafkaConfig {
	cfg := &KafkaConfig{
		Brokers: []string{getEnvOrDefault("KAFKA_BROKERS", "localhost:9092")},
		GroupID: getEnvOrDefault("KAFKA_GROUP_ID", "todo-bot-group"),
	}

	cfg.Topics.TaskCreated = getEnvOrDefault("KAFKA_TOPIC_TASK_CREATED", "task-created")
	cfg.Topics.TaskUpdated = getEnvOrDefault("KAFKA_TOPIC_TASK_UPDATED", "task-updated")
	cfg.Topics.TaskDeleted = getEnvOrDefault("KAFKA_TOPIC_TASK_DELETED", "task-deleted")
	cfg.Topics.TaskOverdue = getEnvOrDefault("KAFKA_TOPIC_TASK_OVERDUE", "task-overdue")
	cfg.Topics.Notifications = getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications")

	return cfg
}
