package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"size:50;default:'pending'"`
	Deadline    time.Time
	Category    string `gorm:"size:100"`
	UserID      int64  `gorm:"not null"`
	Priority    int    `gorm:"default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
