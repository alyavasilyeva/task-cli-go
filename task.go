package main

import (
	"errors"
	"fmt"
	"slices"
	"time"
	"uuid"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
)

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
)

type Task struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

func createTask(tasks []Task, description string) ([]Task, Task, error) {
	if slices.ContainsFunc(tasks, func(task Task) bool {
		return task.Description == description
	}) {
		return tasks, Task{}, fmt.Errorf("%w: %q", ErrTaskAlreadyExists, description)
	}

	now := time.Now().Format(time.RFC3339)
	task := Task{
		ID:          uuid.NewV7().String(),
		Description: description,
		Status:      TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return append(tasks, task), task, nil
}

func deleteTask(tasks []Task, id string) ([]Task, error) {
	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.ID == id
	})
	if index == -1 {
		return tasks, fmt.Errorf("%w: %s", ErrTaskNotFound, id)
	}

	return slices.Delete(tasks, index, index+1), nil
}

func updateTask(tasks []Task, id, description string) ([]Task, error) {
	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.ID == id
	})
	if index == -1 {
		return tasks, fmt.Errorf("%w: %s", ErrTaskNotFound, id)
	}

	tasks[index].Description = description
	tasks[index].UpdatedAt = time.Now().Format(time.RFC3339)

	return tasks, nil
}

func markTask(tasks []Task, id string, status TaskStatus) ([]Task, error) {
	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.ID == id
	})
	if index == -1 {
		return tasks, fmt.Errorf("%w: %s", ErrTaskNotFound, id)
	}

	tasks[index].Status = status
	tasks[index].UpdatedAt = time.Now().Format(time.RFC3339)

	return tasks, nil
}

func filterTasks(tasks []Task, status *TaskStatus) []Task {
	if status == nil {
		return tasks
	}

	filtered := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Status == *status {
			filtered = append(filtered, task)
		}
	}

	return filtered
}
