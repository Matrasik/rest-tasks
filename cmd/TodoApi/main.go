package main

import (
	"Rest-todo/api/handler"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/tasks", handler.GetTasks)
	app.Post("/tasks", handler.PostTasks)
	app.Listen(":8080")
}
