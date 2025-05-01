package handler

import (
	"Rest-todo/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
	"time"
)

type Handler struct {
	Db *pgxpool.Pool
}

func (h *Handler) GetTasks(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json")
	rows, err := h.Db.Query(c.Context(), `SELECT * FROM tasks`)
	defer rows.Close()
	if err != nil {
		log.Printf("failed to get tasks from db: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"tasks": nil,
			"error": "failed to get tasks",
		})
	}
	var tasks []service.Task
	for rows.Next() {
		var t service.Task
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Printf("error scan row from db: %s")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"tasks": nil,
				"error": "failed to get task",
			})
		}
		tasks = append(tasks, t)
	}
	return c.JSON(fiber.Map{
		"tasks": tasks,
		"error": nil,
	})
}

func (h *Handler) PostTasks(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json")
	task := service.Task{}
	err := c.BodyParser(&task)
	if err != nil {
		log.Printf("error parsing json: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"task":  nil,
			"error": "invalid request body",
		})
	}
	task.Status = "new"
	timeNow := time.Now().UTC()
	task.CreatedAt = timeNow
	task.UpdatedAt = timeNow
	var id int
	err = h.Db.QueryRow(c.Context(), `
		INSERT INTO tasks (title, description ,status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id`, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).Scan(&id)
	if err != nil {
		log.Printf("database query insertion error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"task":  nil,
			"error": "failed to create task",
		})
	}
	task.Id = id
	return c.JSON(fiber.Map{
		"task":  task,
		"error": nil,
	})
}

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("invalid id. %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"taskId": nil,
			"error":  "invalid id",
		})
	}
	var deletedId int
	err = h.Db.QueryRow(c.Context(), `DELETE FROM tasks WHERE id = $1 RETURNING ID`, id).Scan(&deletedId)
	if err != nil {
		log.Printf("task with this id not exist. Err:%s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"taskId": nil,
			"error":  "failed to delete task",
		})
	}
	return c.JSON(fiber.Map{
		"taskId": id,
		"error":  nil,
	})
}

func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("invalid id. %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"task":  nil,
			"error": "invalid id",
		})
	}
	task := service.Task{}
	err = c.BodyParser(&task)
	if err != nil {
		log.Printf("error parsing json: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"task":  nil,
			"error": "invalid request body",
		})
	}
	task.UpdatedAt = time.Now().UTC()
	if task.Status != "new" && task.Status != "in_progress" && task.Status != "done" {
		log.Println("invalid task status. Must be new, in_progress or done")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"task":  nil,
			"error": "invalid task status",
		})
	}
	// Если прислали json без description, status или title берем их значения из базы.
	if task.Description == "" || task.Status == "" || task.Title == "" {
		taskExisting := service.Task{}
		err = h.Db.QueryRow(c.Context(), `Select title, description, status From tasks WHERE id = $1`, id).Scan(
			&taskExisting.Title, &taskExisting.Description, &taskExisting.Status,
		)
		if err != nil {
			log.Printf("failed to get task from db: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"task":  nil,
				"error": "failed to get task",
			})
		}
		if task.Description == "" {
			task.Description = taskExisting.Description
		}
		if task.Status == "" {
			task.Status = taskExisting.Status
		}
		if task.Title == "" {
			task.Title = taskExisting.Title
		}
	}
	taskReturn := service.Task{}
	err = h.Db.QueryRow(c.Context(), `
		UPDATE tasks
        SET title = $1, description = $2, status = $3, updated_at = $4
        WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at`,
		task.Title, task.Description, task.Status, task.UpdatedAt, id).Scan(
		&taskReturn.Id, &taskReturn.Title, &taskReturn.Description, &taskReturn.Status,
		&taskReturn.CreatedAt, &taskReturn.UpdatedAt)
	if err != nil {
		log.Printf("database update error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"task":  nil,
			"error": "task not found",
		})
	}
	return c.JSON(fiber.Map{
		"task":  taskReturn,
		"error": nil,
	})
}
