package repository

import (
	"time"
	"to-do-list/internal/models"

	"gorm.io/gorm"
)

type TaskRepositoryInterface interface {
	Create(task *models.Task) error
	GetByID(id uint) (*models.Task, error)
	GetByUserID(userID int64) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
	GetByStatus(userID int64, status string) ([]models.Task, error)
	GetByCategory(userID int64, category string) ([]models.Task, error)
	GetOverdue(userID int64) ([]models.Task, error)
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetByUserID(userID int64) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) GetByStatus(userID int64, status string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND status = ?", userID, status).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetByCategory(userID int64, category string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND category = ?", userID, category).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetOverdue(userID int64) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND deadline < ? AND status != 'completed'", userID, time.Now()).Find(&tasks).Error
	return tasks, err
}
