package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"to-do-list/pkg/systems"
)

const helpText = `Вот список доступных команд:
1. Добавить задание (/add_task)
2. Удалить задание (/delete_task)
3. Список заданий (/show_task)
4. Изменить задание (/change_task)`

func main() {
	token := systems.TakeToken()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Привет!" + " " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName,
			})
		case "/help":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   helpText,
			})
		case "/add_task":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
		case "/delete_task":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
		case "/change_task":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
		case "/show_task":
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
		}
	}
}
