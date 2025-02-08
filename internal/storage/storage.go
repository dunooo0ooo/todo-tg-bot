package storage

import (
	"database/sql"
	"fmt"
	"time"
	"to-do-list/internal/task"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS tasks (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	user_id INTEGER NOT NULL,
    	name TEXT NOT NULL,
    	description TEXT,
    	due_date DATETIME,                    
    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTask(userId int64, taskName string, taskDescription string, dueDate time.Time) (int, error) {
	const op = "storage.sqlite.AddTask"

	res, err := s.db.Exec(`INSERT INTO tasks(user_id, name, description, due_date) VALUES (?, ?, ?, ?)`,
		userId, taskName, taskDescription, dueDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(lastID), nil
}

func (s *Storage) GetTasks(userId int64) ([]task.Task, error) {
	const op = "storage.sqlite.GetTasks"

	rows, err := s.db.Query(`SELECT id, user_id, name, description, due_date, created_at FROM tasks WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.Id, &t.UserId, &t.Name, &t.Description, &t.DueDate, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *Storage) DeleteTask(userId int64, taskName string) error {
	const op = "storage.sqlite.DeleteTask"

	_, err := s.db.Exec(`DELETE FROM tasks WHERE user_id = ? AND name = ?`, userId, taskName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) ChangeTask(taskId int, newName string, newDescription string, newDueDate time.Time) error {
	const op = "storage.sqlite.ChangeTask"

	_, err := s.db.Exec("UPDATE tasks SET name = ?, description = ?, due_date = ? WHERE id = ?", newName, newDescription, newDueDate, taskId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
