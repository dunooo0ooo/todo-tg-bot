package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"to-do-list/internal/bot-service/handlers"
	"to-do-list/pkg/systems"
)

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
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Привет!" + " " + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName,
			})
			if err != nil {
				return
			}
		case "/help":
			handlers.Help(ctx, b, update)
		case "/add_task":
			//_, err := handlers.Add(ctx, b, update)
			//if err != nil {
			//	return
			//}
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Напишите название задачи",
			})
			if err != nil {
				return
			}
		case "/delete_task":
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
			if err != nil {
				return
			}
		case "/change_task":
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
			if err != nil {
				return
			}
		case "/show_task":
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "123",
			})
			if err != nil {
				return
			}
		default:
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Используйте команду /help, чтобы узнать возможности бота",
			})
		}

	}
}
