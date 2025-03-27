package service

import (
	"context"
	"errors"
	"log"
	"time"
	"to-do-list/internal/models"
	"to-do-list/internal/repository"
	"to-do-list/pkg/kafka"
)

type TaskService struct {
	repo     repository.TaskRepositoryInterface
	producer kafka.ProducerInterface
}

func NewTaskService(repo repository.TaskRepositoryInterface, producer kafka.ProducerInterface) *TaskService {
	return &TaskService{
		repo:     repo,
		producer: producer,
	}
}

func (s *TaskService) CreateTask(title string, description string, userID int64, deadline time.Time, category string) (*models.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		UserID:      userID,
		Deadline:    deadline,
		Category:    category,
		Status:      "pending",
		Priority:    1,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	if err := s.producer.SendTaskCreated(context.Background(), task); err != nil {
		log.Printf("Error sending task created event: %v", err)
	}

	return task, nil
}

func (s *TaskService) GetTask(id uint) (*models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) GetUserTasks(userID int64) ([]models.Task, error) {
	return s.repo.GetByUserID(userID)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}

	if err := s.repo.Update(task); err != nil {
		return err
	}

	if err := s.producer.SendTaskUpdated(context.Background(), task); err != nil {
		log.Printf("Error sending task updated event: %v", err)
	}

	return nil
}

func (s *TaskService) DeleteTask(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if err := s.producer.SendTaskDeleted(context.Background(), id); err != nil {
		log.Printf("Error sending task deleted event: %v", err)
	}

	return nil
}

func (s *TaskService) GetTasksByStatus(userID int64, status string) ([]models.Task, error) {
	return s.repo.GetByStatus(userID, status)
}

func (s *TaskService) GetTasksByCategory(userID int64, category string) ([]models.Task, error) {
	return s.repo.GetByCategory(userID, category)
}

func (s *TaskService) GetOverdueTasks(userID int64) ([]models.Task, error) {
	return s.repo.GetOverdue(userID)
}

func (s *TaskService) UpdateTaskStatus(id uint, status string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	task.Status = status
	if err := s.repo.Update(task); err != nil {
		return err
	}

	if err := s.producer.SendTaskUpdated(context.Background(), task); err != nil {
		log.Printf("Error sending task updated event: %v", err)
	}

	return nil
}

func (s *TaskService) UpdateTaskPriority(id uint, priority int) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	task.Priority = priority
	if err := s.repo.Update(task); err != nil {
		return err
	}

	if err := s.producer.SendTaskUpdated(context.Background(), task); err != nil {
		log.Printf("Error sending task updated event: %v", err)
	}

	return nil
}
