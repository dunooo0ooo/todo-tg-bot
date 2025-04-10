package service

import (
	"context"
	"testing"
	"time"
	"to-do-list/internal/models"
)

type MockProducer struct{}

func (m *MockProducer) SendTaskCreated(ctx context.Context, task *models.Task) error { return nil }
func (m *MockProducer) SendTaskUpdated(ctx context.Context, task *models.Task) error { return nil }
func (m *MockProducer) SendTaskDeleted(ctx context.Context, taskID uint) error       { return nil }
func (m *MockProducer) SendTaskOverdue(ctx context.Context, task *models.Task) error { return nil }
func (m *MockProducer) SendNotification(ctx context.Context, userID int64, message string) error {
	return nil
}
func (m *MockProducer) Close() error { return nil }

// MockTaskRepository - мок репозитория для тестирования
type MockTaskRepository struct {
	tasks map[uint]*models.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make(map[uint]*models.Task),
	}
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	task.ID = uint(len(m.tasks) + 1)
	m.tasks[task.ID] = task
	return nil
}

func (m *MockTaskRepository) GetByID(id uint) (*models.Task, error) {
	if task, exists := m.tasks[id]; exists {
		return task, nil
	}
	return nil, nil
}

func (m *MockTaskRepository) GetByUserID(userID int64) ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range m.tasks {
		if task.UserID == userID {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	if _, exists := m.tasks[task.ID]; exists {
		m.tasks[task.ID] = task
		return nil
	}
	return nil
}

func (m *MockTaskRepository) Delete(id uint) error {
	delete(m.tasks, id)
	return nil
}

func (m *MockTaskRepository) GetByStatus(userID int64, status string) ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range m.tasks {
		if task.UserID == userID && task.Status == status {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) GetByCategory(userID int64, category string) ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range m.tasks {
		if task.UserID == userID && task.Category == category {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) GetOverdue(userID int64) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	for _, task := range m.tasks {
		if task.UserID == userID && task.Deadline.Before(now) && task.Status != "completed" {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	mockProducer := &MockProducer{}
	service := NewTaskService(mockRepo, mockProducer)

	tests := []struct {
		name        string
		title       string
		description string
		userID      int64
		deadline    time.Time
		category    string
		wantErr     bool
	}{
		{
			name:        "Valid task",
			title:       "Test Task",
			description: "Test Description",
			userID:      123,
			deadline:    time.Now().Add(24 * time.Hour),
			category:    "test",
			wantErr:     false,
		},
		{
			name:        "Empty title",
			title:       "",
			description: "Test Description",
			userID:      123,
			deadline:    time.Now().Add(24 * time.Hour),
			category:    "test",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.CreateTask(tt.title, tt.description, tt.userID, tt.deadline, tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && task == nil {
				t.Error("CreateTask() returned nil task without error")
			}
		})
	}
}

func TestTaskService_GetUserTasks(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	mockProducer := &MockProducer{}
	service := NewTaskService(mockRepo, mockProducer)

	// Создаем тестовые задачи
	userID := int64(123)
	task1 := &models.Task{
		Title:       "Task 1",
		Description: "Description 1",
		UserID:      userID,
		Deadline:    time.Now().Add(24 * time.Hour),
		Category:    "test",
	}
	task2 := &models.Task{
		Title:       "Task 2",
		Description: "Description 2",
		UserID:      userID,
		Deadline:    time.Now().Add(48 * time.Hour),
		Category:    "test",
	}

	mockRepo.Create(task1)
	mockRepo.Create(task2)

	tasks, err := service.GetUserTasks(userID)
	if err != nil {
		t.Errorf("GetUserTasks() error = %v", err)
		return
	}

	if len(tasks) != 2 {
		t.Errorf("GetUserTasks() returned %d tasks, want 2", len(tasks))
	}
}

func TestTaskService_UpdateTaskStatus(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	mockProducer := &MockProducer{}
	service := NewTaskService(mockRepo, mockProducer)

	// Создаем тестовую задачу
	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		UserID:      123,
		Deadline:    time.Now().Add(24 * time.Hour),
		Category:    "test",
		Status:      "pending",
	}

	mockRepo.Create(task)

	err := service.UpdateTaskStatus(task.ID, "completed")
	if err != nil {
		t.Errorf("UpdateTaskStatus() error = %v", err)
		return
	}

	updatedTask, err := mockRepo.GetByID(task.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
		return
	}

	if updatedTask.Status != "completed" {
		t.Errorf("UpdateTaskStatus() failed to update status, got %v, want completed", updatedTask.Status)
	}
}
