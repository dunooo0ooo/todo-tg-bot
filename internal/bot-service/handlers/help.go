package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const helpText = `Вот список доступных команд:
1. Добавить задание (/add_task)
2. Удалить задание (/delete_task)
3. Список заданий (/show_task)
4. Изменить задание (/change_task)`

func Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
}
