package task

import "time"

type Task struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      int64     `gorm:"index" json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
