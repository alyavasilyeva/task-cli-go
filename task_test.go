package main

import (
	"errors"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	tasks, created, err := createTask(nil, "Buy milk")
	if err != nil {
		t.Fatalf("createTask() error = %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("createTask() returned %d tasks, want 1", len(tasks))
	}
	if created.ID == "" {
		t.Error("createTask() returned an empty ID")
	}
	if created.Description != "Buy milk" {
		t.Errorf("createTask() description = %q, want %q", created.Description, "Buy milk")
	}
	if created.Status != TaskStatusTodo {
		t.Errorf("createTask() status = %q, want %q", created.Status, TaskStatusTodo)
	}
	if created.CreatedAt != created.UpdatedAt {
		t.Errorf("createTask() timestamps differ: created_at=%q updated_at=%q", created.CreatedAt, created.UpdatedAt)
	}
	if _, err := time.Parse(time.RFC3339, created.CreatedAt); err != nil {
		t.Errorf("createTask() created_at is not RFC3339: %v", err)
	}
}

func TestCreateTaskRejectsDuplicateDescription(t *testing.T) {
	existing := []Task{{ID: "1", Description: "Buy milk"}}

	tasks, _, err := createTask(existing, "Buy milk")

	if !errors.Is(err, ErrTaskAlreadyExists) {
		t.Fatalf("createTask() error = %v, want ErrTaskAlreadyExists", err)
	}
	if len(tasks) != len(existing) {
		t.Errorf("createTask() returned %d tasks, want %d", len(tasks), len(existing))
	}
}

func TestUpdateTask(t *testing.T) {
	tasks := []Task{{
		ID:          "1",
		Description: "Old description",
		UpdatedAt:   "2026-01-01T00:00:00Z",
	}}

	updated, err := updateTask(tasks, "1", "New description")

	if err != nil {
		t.Fatalf("updateTask() error = %v", err)
	}
	if updated[0].Description != "New description" {
		t.Errorf("updateTask() description = %q, want %q", updated[0].Description, "New description")
	}
	if updated[0].UpdatedAt == "2026-01-01T00:00:00Z" {
		t.Error("updateTask() did not update updated_at")
	}
}

func TestMarkTask(t *testing.T) {
	tasks := []Task{{ID: "1", Status: TaskStatusTodo}}

	updated, err := markTask(tasks, "1", TaskStatusDone)

	if err != nil {
		t.Fatalf("markTask() error = %v", err)
	}
	if updated[0].Status != TaskStatusDone {
		t.Errorf("markTask() status = %q, want %q", updated[0].Status, TaskStatusDone)
	}
}

func TestDeleteTask(t *testing.T) {
	tasks := []Task{{ID: "1"}, {ID: "2"}}

	updated, err := deleteTask(tasks, "1")

	if err != nil {
		t.Fatalf("deleteTask() error = %v", err)
	}
	if len(updated) != 1 || updated[0].ID != "2" {
		t.Errorf("deleteTask() returned %#v, want only task 2", updated)
	}
}

func TestTaskOperationsReturnNotFound(t *testing.T) {
	tests := map[string]func() error{
		"delete": func() error {
			_, err := deleteTask(nil, "missing")
			return err
		},
		"update": func() error {
			_, err := updateTask(nil, "missing", "description")
			return err
		},
		"mark": func() error {
			_, err := markTask(nil, "missing", TaskStatusDone)
			return err
		},
	}

	for name, operation := range tests {
		t.Run(name, func(t *testing.T) {
			if err := operation(); !errors.Is(err, ErrTaskNotFound) {
				t.Errorf("operation error = %v, want ErrTaskNotFound", err)
			}
		})
	}
}

func TestFilterTasks(t *testing.T) {
	tasks := []Task{
		{ID: "1", Status: TaskStatusTodo},
		{ID: "2", Status: TaskStatusDone},
		{ID: "3", Status: TaskStatusTodo},
	}
	status := TaskStatusTodo

	filtered := filterTasks(tasks, &status)

	if len(filtered) != 2 || filtered[0].ID != "1" || filtered[1].ID != "3" {
		t.Errorf("filterTasks() returned %#v, want todo tasks 1 and 3", filtered)
	}
}
