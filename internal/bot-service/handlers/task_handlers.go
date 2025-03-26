package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"to-do-list/internal/service"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TaskHandlers struct {
	taskService *service.TaskService
	stateMgr    *StateManager
}

func NewTaskHandlers(taskService *service.TaskService) *TaskHandlers {
	return &TaskHandlers{
		taskService: taskService,
		stateMgr:    NewStateManager(),
	}
}

func (h *TaskHandlers) AddTask(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := update.Message.Chat.ID
	state := &UserState{
		CurrentStep: StepAddingTitle,
	}
	h.stateMgr.SetState(userID, state)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   "Введите название задачи:",
	})
	return err
}

func (h *TaskHandlers) ShowTasks(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := update.Message.Chat.ID
	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "У вас пока нет задач. Используйте /add_task для создания новой задачи.",
		})
		return err
	}

	var message strings.Builder
	message.WriteString("Ваши задачи:\n\n")
	for i, task := range tasks {
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, task.Title))
		message.WriteString(fmt.Sprintf("   Описание: %s\n", task.Description))
		message.WriteString(fmt.Sprintf("   Дедлайн: %s\n", task.Deadline.Format("02.01.2006 15:04")))
		message.WriteString(fmt.Sprintf("   Категория: %s\n", task.Category))
		message.WriteString(fmt.Sprintf("   Статус: %s\n\n", task.Status))
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   message.String(),
	})
	return err
}

func (h *TaskHandlers) DeleteTask(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := update.Message.Chat.ID
	state := &UserState{
		CurrentStep: StepDeletingTask,
	}
	h.stateMgr.SetState(userID, state)

	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "У вас нет задач для удаления.",
		})
		return err
	}

	var message strings.Builder
	message.WriteString("Выберите номер задачи для удаления:\n\n")
	for i, task := range tasks {
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, task.Title))
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   message.String(),
	})
	return err
}

func (h *TaskHandlers) UpdateTaskStatus(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := update.Message.Chat.ID
	state := &UserState{
		CurrentStep: StepEditingTask,
	}
	h.stateMgr.SetState(userID, state)

	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "У вас нет задач для обновления.",
		})
		return err
	}

	var message strings.Builder
	message.WriteString("Выберите номер задачи для обновления статуса:\n\n")
	for i, task := range tasks {
		message.WriteString(fmt.Sprintf("%d. %s (текущий статус: %s)\n", i+1, task.Title, task.Status))
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   message.String(),
	})
	return err
}

func (h *TaskHandlers) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) error {
	userID := update.Message.Chat.ID
	text := update.Message.Text
	state := h.stateMgr.GetState(userID)

	if state == nil {
		state = &UserState{CurrentStep: StepIdle}
	}

	switch state.CurrentStep {
	case StepAddingTitle:
		state.TaskTitle = text
		state.CurrentStep = StepAddingDesc
		h.stateMgr.SetState(userID, state)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Введите описание задачи:",
		})
		return err

	case StepAddingDesc:
		state.TaskDesc = text
		state.CurrentStep = StepAddingDeadline
		h.stateMgr.SetState(userID, state)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Введите дедлайн задачи в формате ДД.ММ.ГГГГ ЧЧ:ММ:",
		})
		return err

	case StepAddingDeadline:
		deadline, err := time.Parse("02.01.2006 15:04", text)
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Неверный формат даты. Попробуйте еще раз в формате ДД.ММ.ГГГГ ЧЧ:ММ:",
			})
			return err
		}
		state.TaskDeadline = deadline
		state.CurrentStep = StepAddingCategory
		h.stateMgr.SetState(userID, state)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Введите категорию задачи:",
		})
		return err

	case StepAddingCategory:
		state.TaskCategory = text
		task, err := h.taskService.CreateTask(
			state.TaskTitle,
			state.TaskDesc,
			userID,
			state.TaskDeadline,
			state.TaskCategory,
		)
		if err != nil {
			return err
		}
		h.stateMgr.DeleteState(userID)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   fmt.Sprintf("Задача успешно создана!\nНазвание: %s\nДедлайн: %s", task.Title, task.Deadline.Format("02.01.2006 15:04")),
		})
		return err

	case StepDeletingTask:
		taskNum, err := strconv.Atoi(text)
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Пожалуйста, введите номер задачи:",
			})
			return err
		}
		tasks, err := h.taskService.GetUserTasks(userID)
		if err != nil {
			return err
		}
		if taskNum < 1 || taskNum > len(tasks) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Неверный номер задачи. Попробуйте еще раз:",
			})
			return err
		}
		task := tasks[taskNum-1]
		err = h.taskService.DeleteTask(task.ID)
		if err != nil {
			return err
		}
		h.stateMgr.DeleteState(userID)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Задача успешно удалена!",
		})
		return err

	case StepEditingTask:
		taskNum, err := strconv.Atoi(text)
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Пожалуйста, введите номер задачи:",
			})
			return err
		}
		tasks, err := h.taskService.GetUserTasks(userID)
		if err != nil {
			return err
		}
		if taskNum < 1 || taskNum > len(tasks) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Неверный номер задачи. Попробуйте еще раз:",
			})
			return err
		}
		task := tasks[taskNum-1]
		state.EditingTaskID = task.ID
		state.CurrentStep = "updating_status"
		h.stateMgr.SetState(userID, state)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Введите новый статус задачи (pending/in_progress/completed):",
		})
		return err

	case "updating_status":
		status := text
		if status != "pending" && status != "in_progress" && status != "completed" {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userID,
				Text:   "Неверный статус. Введите один из: pending, in_progress, completed",
			})
			return err
		}
		err := h.taskService.UpdateTaskStatus(state.EditingTaskID, status)
		if err != nil {
			return err
		}
		h.stateMgr.DeleteState(userID)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   "Статус задачи успешно обновлен!",
		})
		return err
	}

	return nil
}
