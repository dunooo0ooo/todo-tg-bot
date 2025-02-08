package handlers

//
//import (
//	"context"
//	"fmt"
//	"github.com/go-telegram/bot"
//	"github.com/go-telegram/bot/models"
//	"time"
//)
//
//
//var userState = make(map[int64]string)
//var userTasks = make(map[int64]*InputTask)
//
//type InputTask struct {
//	TaskName        string
//	TaskDescription string
//	DueDate         time.Time
//}
//
//
//func Add(ctx context.Context, b *bot.Bot, update *models.Update) (int, error) {
//	userID := update.Message.Chat.ID
//
//	state, exists := userState[userID]
//
//	if !exists || state == "" {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: userID,
//			Text:   "Напишите название задачи",
//		})
//
//		userState[userID] = "waiting_for_task_name"
//		userTasks[userID] = &InputTask{}
//		return 0, nil
//	}
//
//	task := userTasks[userID]
//
//	if state == "waiting_for_task_name" {
//		task.TaskName = update.Message.Text
//
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: userID,
//			Text:   "Теперь напишите описание задачи",
//		})
//
//		userState[userID] = "waiting_for_task_description"
//		return 0, nil
//	}
//
//	if state == "waiting_for_task_description" {
//		task.TaskDescription = update.Message.Text
//
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: userID,
//			Text:   "Укажите дату выполнения задачи в формате ГГГГ-ММ-ДД",
//		})
//
//		userState[userID] = "waiting_for_due_date"
//		return 0, nil
//	}
//
//	if state == "waiting_for_due_date" {
//		dueDateStr := update.Message.Text
//
//		dueDate, err := time.Parse("2006-01-02", dueDateStr)
//		if err != nil {
//			b.SendMessage(ctx, &bot.SendMessageParams{
//				ChatID: userID,
//				Text:   "Некорректный формат даты. Попробуйте ещё раз: ГГГГ-ММ-ДД",
//			})
//			return 0, nil
//		}
//
//		task.DueDate = dueDate
//
//		const op = "storage.sqlite.AddTask"
//
//		stmt, err := db.Prepare(`insert into tasks(user_id ,name, description, due_date) values(?, ?, ?, ?)`)
//
//		if err != nil {
//			return 0, fmt.Errorf("%s: %w", op, err)
//		}
//
//		res, err := stmt.Exec(userID, task.TaskName, task.TaskDescription, task.DueDate)
//		if err != nil {
//			return 0, fmt.Errorf("%s: %w", op, err)
//		}
//
//		lastID, err := res.LastInsertId()
//		if err != nil {
//			return 0, fmt.Errorf("%s: %w", op, err)
//		}
//
//		return int(lastID), nil
//
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID: userID,
//			Text:   "Задача успешно добавлена!",
//		})
//
//		delete(userState, userID)
//		delete(userTasks, userID)
//	}
//	return 0, nil
//}
