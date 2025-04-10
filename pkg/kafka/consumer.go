package kafka

import (
	"context"
	"encoding/json"
	"log"
	"to-do-list/internal/models"
	"to-do-list/pkg/config"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	config   *config.KafkaConfig
	handlers map[string]func([]byte) error
}

func NewConsumer(cfg *config.KafkaConfig) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, err
	}

	c := &Consumer{
		consumer: consumer,
		config:   cfg,
		handlers: make(map[string]func([]byte) error),
	}

	// Регистрируем обработчики
	c.handlers[cfg.Topics.TaskCreated] = c.handleTaskCreated
	c.handlers[cfg.Topics.TaskUpdated] = c.handleTaskUpdated
	c.handlers[cfg.Topics.TaskDeleted] = c.handleTaskDeleted
	c.handlers[cfg.Topics.TaskOverdue] = c.handleTaskOverdue
	c.handlers[cfg.Topics.Notifications] = c.handleNotification

	return c, nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) Start(ctx context.Context) error {
	topics := []string{
		c.config.Topics.TaskCreated,
		c.config.Topics.TaskUpdated,
		c.config.Topics.TaskDeleted,
		c.config.Topics.TaskOverdue,
		c.config.Topics.Notifications,
	}

	for {
		err := c.consumer.Consume(ctx, topics, c)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		handler, exists := c.handlers[message.Topic]
		if !exists {
			log.Printf("No handler for topic: %s", message.Topic)
			continue
		}

		if err := handler(message.Value); err != nil {
			log.Printf("Error handling message from topic %s: %v", message.Topic, err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func (c *Consumer) handleTaskCreated(data []byte) error {
	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}
	log.Printf("Task created: %+v", task)
	return nil
}

func (c *Consumer) handleTaskUpdated(data []byte) error {
	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}
	log.Printf("Task updated: %+v", task)
	return nil
}

func (c *Consumer) handleTaskDeleted(data []byte) error {
	var msg struct {
		TaskID uint `json:"task_id"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	log.Printf("Task deleted: %d", msg.TaskID)
	return nil
}

func (c *Consumer) handleTaskOverdue(data []byte) error {
	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}
	log.Printf("Task overdue: %+v", task)
	return nil
}

func (c *Consumer) handleNotification(data []byte) error {
	var notification struct {
		UserID  int64  `json:"user_id"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &notification); err != nil {
		return err
	}
	log.Printf("Notification for user %d: %s", notification.UserID, notification.Message)
	return nil
}
