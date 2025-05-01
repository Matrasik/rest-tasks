package handler

import (
	"Rest-todo/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Handler struct {
	Db *pgxpool.Pool
}

func (h *Handler) GetTasks(c *fiber.Ctx) error {
	return c.SendString("get method")
}

func (h *Handler) PostTasks(c *fiber.Ctx) error {
	//TODO проверка на существование в базе
	c.Set("Content-type", "application/json")
	task := service.Task{}
	err := c.BodyParser(&task)
	if err != nil {
		log.Printf("error parsing json: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	task.Status = "new"
	time := time.Now().UTC()
	task.CreatedAt = time
	task.UpdatedAt = time
	var id int
	err = h.Db.QueryRow(c.Context(), `
		INSERT INTO tasks (title, description ,status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id`, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).Scan(&id)
	if err != nil {
		log.Printf("database query insertion error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create task",
		})
	}
	return c.JSON(fiber.Map{
		"json": task,
	})
}
