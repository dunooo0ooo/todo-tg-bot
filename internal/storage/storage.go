package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)

	stmt, err := db.Prepare(`
    CREATE TABLE tasks (
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

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTask(userId int64, taskName string, taskDescription string, dueDate time.Time) (int, error) {
	const op = "storage.sqlite.AddTask"

	stmt, err := s.db.Prepare(`insert into tasks(user_id ,name, description, due_date) values(?, ?, ?, ?)`)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(userId, taskName, taskDescription, dueDate)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(lastID), nil
}
