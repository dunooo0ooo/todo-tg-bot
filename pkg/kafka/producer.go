package kafka

import (
	"context"
	"encoding/json"
	"to-do-list/internal/models"
	"to-do-list/pkg/config"

	"github.com/IBM/sarama"
)

type ProducerInterface interface {
	SendTaskCreated(ctx context.Context, task *models.Task) error
	SendTaskUpdated(ctx context.Context, task *models.Task) error
	SendTaskDeleted(ctx context.Context, taskID uint) error
	SendTaskOverdue(ctx context.Context, task *models.Task) error
	SendNotification(ctx context.Context, userID int64, message string) error
	Close() error
}

type Producer struct {
	producer sarama.SyncProducer
	config   *config.KafkaConfig
}

func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		config:   cfg,
	}, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) SendTaskCreated(ctx context.Context, task *models.Task) error {
	msg, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.config.Topics.TaskCreated,
		Value: sarama.StringEncoder(msg),
	})
	return err
}

func (p *Producer) SendTaskUpdated(ctx context.Context, task *models.Task) error {
	msg, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.config.Topics.TaskUpdated,
		Value: sarama.StringEncoder(msg),
	})
	return err
}

func (p *Producer) SendTaskDeleted(ctx context.Context, taskID uint) error {
	msg, err := json.Marshal(map[string]uint{"task_id": taskID})
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.config.Topics.TaskDeleted,
		Value: sarama.StringEncoder(msg),
	})
	return err
}

func (p *Producer) SendTaskOverdue(ctx context.Context, task *models.Task) error {
	msg, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.config.Topics.TaskOverdue,
		Value: sarama.StringEncoder(msg),
	})
	return err
}

func (p *Producer) SendNotification(ctx context.Context, userID int64, message string) error {
	notification := struct {
		UserID  int64  `json:"user_id"`
		Message string `json:"message"`
	}{
		UserID:  userID,
		Message: message,
	}

	msg, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.config.Topics.Notifications,
		Value: sarama.StringEncoder(msg),
	})
	return err
}
