package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"to-do-list/internal/storage"
	"to-do-list/internal/todo-service/handlers/add"
	"to-do-list/internal/todo-service/handlers/change"
	"to-do-list/internal/todo-service/handlers/delete"
	"to-do-list/internal/todo-service/handlers/show"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := storage.New(dsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Post("/add/task", add.Add(db))
	router.Get("/tasks", show.ShowTasks(db))
	router.Delete("/delete/task", delete.DeleteHandler(db))
	router.Put("/change/task", change.ChangeHandler(db))

	srv := &http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	log.Println("API server is running on http://localhost:3000")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}
