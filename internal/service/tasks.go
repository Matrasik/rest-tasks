package service

import "time"

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"-"`
	CratedAt    time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
