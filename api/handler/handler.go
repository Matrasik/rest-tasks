package handler

import (
	"Rest-todo/internal/service"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"time"
)

func GetTasks(c *fiber.Ctx) error {
	return c.SendString("get method")
}

func PostTasks(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json")
	Task := service.Task{}
	err := c.BodyParser(&Task)
	if err != nil {
		log.Printf("error parsing json: %s", err)
		c.Status(http.StatusInternalServerError).SendString("InternalServerError")
	}
	Task.Status = "new"
	time := time.Now()
	Task.CratedAt = time
	Task.UpdatedAt = time
	return c.JSONP(fiber.Map{
		"json": Task,
	})
}
