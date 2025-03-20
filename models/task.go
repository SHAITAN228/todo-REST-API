package models

import (
	"github.com/go-playground/validator/v10"
)

type Task struct {
	ID        int64  `json:"id"`
	Title     string `json:"title" binding:"required" validate:"required,min=2,max=100"`
	Content   string `json:"content" validate:"min=0,max=200"`
	Completed bool   `json:"completed"`
}

func (task *Task) Validate() error {
	return validator.New().Struct(task)
}
