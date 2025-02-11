package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"to-do-list/internal/storage"
	"to-do-list/internal/todo-service/handlers"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=todo_list port=5432 sslmode=disable TimeZone=UTC"

	db, err := storage.New(dsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Post("/add/task", handlers.Add(db))
	router.Get("/tasks", handlers.ShowTasks(db))
	router.Delete("/delete/task", handlers.DeleteHandler(db))
	router.Put("/change/task", handlers.ChangeHandler(db))

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
