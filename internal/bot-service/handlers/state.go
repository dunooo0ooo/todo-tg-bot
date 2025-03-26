package handlers

import (
	"sync"
	"time"
)

type UserState struct {
	CurrentStep   string
	TaskTitle     string
	TaskDesc      string
	TaskDeadline  time.Time
	TaskCategory  string
	LastMessageID int
	EditingTaskID uint
}

type StateManager struct {
	states map[int64]*UserState
	mu     sync.RWMutex
}

func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[int64]*UserState),
	}
}

func (sm *StateManager) GetState(userID int64) *UserState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.states[userID]
}

func (sm *StateManager) SetState(userID int64, state *UserState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[userID] = state
}

func (sm *StateManager) DeleteState(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, userID)
}

const (
	StepIdle           = "idle"
	StepAddingTitle    = "adding_title"
	StepAddingDesc     = "adding_desc"
	StepAddingDeadline = "adding_deadline"
	StepAddingCategory = "adding_category"
	StepEditingTask    = "editing_task"
	StepDeletingTask   = "deleting_task"
)
