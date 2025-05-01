package main

import (
	"Rest-todo/api/handler"
	"Rest-todo/internal/db"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSL")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
	dbPool, err := db.CreateConnPool(context.Background(), dsn)
	defer dbPool.Close()
	if err != nil {
		log.Fatalf("failed create connection pool. Error: %s", err)
	}
	m, err := migrate.New(
		"file://internal/db/migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to initialize migration. Err: %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations. Err: %s", err)
	}
	log.Printf("migration sucess")
	h := handler.Handler{
		Db: dbPool,
	}

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${time}][${ip}] ${method} ${path}\n",
	}))
	app.Get("/tasks", h.GetTasks)
	app.Post("/tasks", h.PostTasks)
	app.Delete("/tasks/:id", h.DeleteTask)
	app.Put("/tasks/:id", h.UpdateTask)
	err = app.Listen(":8080")
	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}
	log.Println("Server started at :8080")
}
