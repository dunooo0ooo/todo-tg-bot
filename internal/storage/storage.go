package storage

import (
	"fmt"
	"time"
	"to-do-list/internal/task"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.gorm.NewStorage"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.AutoMigrate(&task.Task{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTask(userId int64, name, description string, dueDate time.Time) (int64, error) {
	const op = "storage.gorm.AddTask"

	task := task.Task{
		UserID:      userId,
		Name:        name,
		Description: description,
		DueDate:     dueDate,
	}

	result := s.db.Create(&task)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return task.ID, nil
}

func (s *Storage) GetTasks(userId int64) ([]task.Task, error) {
	const op = "storage.gorm.GetTasks"

	var tasks []task.Task
	result := s.db.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return tasks, nil
}

func (s *Storage) DeleteTask(userId int64, taskName string) error {
	const op = "storage.gorm.DeleteTask"

	result := s.db.Delete(&task.Task{}, userId, taskName)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (s *Storage) ChangeTask(taskId int64, newName, newDescription string, newDueDate time.Time) error {
	const op = "storage.gorm.ChangeTask"

	result := s.db.Model(&task.Task{}).
		Where("id = ?", taskId).
		Updates(task.Task{
			Name:        newName,
			Description: newDescription,
			DueDate:     newDueDate,
		})

	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}
