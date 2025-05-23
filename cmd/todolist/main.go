package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"to-do-list/internal/bot-service/handlers"
	taskmodel "to-do-list/internal/models"
	"to-do-list/internal/repository"
	"to-do-list/internal/service"
	"to-do-list/pkg/config"
	"to-do-list/pkg/kafka"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var taskHandlers *handlers.TaskHandlers

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := conf.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&taskmodel.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	kafkaConfig := config.LoadKafkaConfig()

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, producer)

	taskHandlers = handlers.NewTaskHandlers(taskService)

	token := conf.BotToken

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.Chat.ID
	text := update.Message.Text

	switch text {
	case "/start":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   fmt.Sprintf("Привет! Я бот для управления задачами. Используйте /help для просмотра доступных команд."),
		})
		if err != nil {
			log.Printf("Error sending start message: %v", err)
		}
	case "/help":
		helpText := `Доступные команды:
/add_task - Добавить новую задачу
/show_tasks - Показать список задач
/delete_task - Удалить задачу
/update_status - Изменить статус задачи
/help - Показать это сообщение`
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   helpText,
		})
		if err != nil {
			log.Printf("Error sending help message: %v", err)
		}

	case "/add_task":
		err := taskHandlers.AddTask(ctx, b, update)
		if err != nil {
			log.Printf("Error handling add task: %v", err)
		}

	case "/show_tasks":
		err := taskHandlers.ShowTasks(ctx, b, update)
		if err != nil {
			log.Printf("Error handling show tasks: %v", err)
		}

	case "/delete_task":
		err := taskHandlers.DeleteTask(ctx, b, update)
		if err != nil {
			log.Printf("Error handling delete task: %v", err)
		}

	case "/update_status":
		err := taskHandlers.UpdateTaskStatus(ctx, b, update)
		if err != nil {
			log.Printf("Error handling update status: %v", err)
		}

	default:
		err := taskHandlers.HandleMessage(ctx, b, update)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}
