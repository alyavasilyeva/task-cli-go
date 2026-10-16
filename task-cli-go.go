package main

import (
	"encoding/json/v2"
	"errors"
	"flag"
	"fmt"
	"os"
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

type Task struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

func loadTasks(filename string) ([]Task, error) {
	byteValue, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}
	if len(byteValue) == 0 {
		return nil, nil
	}

	var tasks []Task
	if err := json.Unmarshal(byteValue, &tasks); err != nil {
		return nil, fmt.Errorf("decode %s: %w", filename, err)
	}
	return tasks, nil
}

func create(description string, filename string) error {
	tasks, err := loadTasks(filename)
	if err != nil {
		return err
	}

	// check if task already exists
	index := slices.IndexFunc(tasks, func(t Task) bool { return t.Description == description })
	if index != -1 {
		fmt.Println("Task already exists")
		return nil
	}

	// create task
	task := Task{
		Id:          uuid.NewV7().String(),
		Description: description,
		Status:      TaskStatusTodo,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	tasks = append(tasks, task)

	// marshal tasks into json
	byteValue, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	// write tasks to file
	if err := os.WriteFile(filename, byteValue, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}

	fmt.Println("Task created:", task.Id)
	return nil
}

func delete(id string, filename string) error {
	tasks, err := loadTasks(filename)
	if err != nil {
		return err
	}

	// check if task exists
	index := slices.IndexFunc(tasks, func(t Task) bool { return t.Id == id })
	if index == -1 {
		fmt.Println("Task not found")
		return nil
	}

	// delete task
	tasks = slices.Delete(tasks, index, index+1)

	// marshal tasks into json
	byteValue, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	// write tasks to file
	if err := os.WriteFile(filename, byteValue, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}

	fmt.Println("Task deleted:", id)
	return nil
}

func update(id string, description string, filename string) error {
	tasks, err := loadTasks(filename)
	if err != nil {
		return err
	}

	// check if task exists
	index := slices.IndexFunc(tasks, func(t Task) bool { return t.Id == id })
	if index == -1 {
		fmt.Println("Task not found")
		return nil
	}

	// update task
	tasks[index].Description = description
	tasks[index].UpdatedAt = time.Now().Format(time.RFC3339)

	// marshal tasks into json
	byteValue, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	// write tasks to file
	if err := os.WriteFile(filename, byteValue, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}

	fmt.Println("Task updated:", id, description)
	return nil
}

func mark(id string, status TaskStatus, filename string) error {
	tasks, err := loadTasks(filename)
	if err != nil {
		return err
	}

	// check if task exists
	index := slices.IndexFunc(tasks, func(t Task) bool { return t.Id == id })
	if index == -1 {
		fmt.Println("Task not found")
		return nil
	}

	// update task
	tasks[index].Status = status
	tasks[index].UpdatedAt = time.Now().Format(time.RFC3339)

	// marshal tasks into json
	byteValue, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	// write tasks to file
	if err := os.WriteFile(filename, byteValue, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}

	fmt.Println("Task marked:", id, status)
	return nil
}

func list(filename string) error {
	tasks, err := loadTasks(filename)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	var taskErrors []error
	fmt.Printf("|%-36s|%-20s|%-12s|%-20s|%-20s|\n", "ID", "Description", "Status", "CreatedAt", "UpdatedAt")
	for _, task := range tasks {
		createdAt, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil {
			taskErrors = append(taskErrors, fmt.Errorf("task %s: invalid created_at %q: %w", task.Id, task.CreatedAt, err))
			continue
		}
		updatedAt, err := time.Parse(time.RFC3339, task.UpdatedAt)
		if err != nil {
			taskErrors = append(taskErrors, fmt.Errorf("task %s: invalid updated_at %q: %w", task.Id, task.UpdatedAt, err))
			continue
		}

		fmt.Printf("|%-36s|%-20s|%-12s|%-20s|%-20s|\n", task.Id, task.Description, task.Status, createdAt.Format(time.DateTime), updatedAt.Format(time.DateTime))
	}
	return errors.Join(taskErrors...)
}

func main() {
	const filename = "tasks.json"

	createFlag := flag.NewFlagSet("create", flag.ExitOnError)
	createFlag.Usage = func() {
		fmt.Fprintln(createFlag.Output(), "Usage: task-cli-go create <description>")
	}

	listFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listFlag.Usage = func() {
		fmt.Fprintln(listFlag.Output(), "Usage: task-cli-go list")
	}

	deleteFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteFlag.Usage = func() {
		fmt.Fprintln(deleteFlag.Output(), "Usage: task-cli-go delete <id>")
	}

	updateFlag := flag.NewFlagSet("update", flag.ExitOnError)
	updateFlag.Usage = func() {
		fmt.Fprintln(updateFlag.Output(), "Usage: task-cli-go update <id> <description>")
	}

	markFlag := flag.NewFlagSet("mark", flag.ExitOnError)
	markFlag.Usage = func() {
		fmt.Fprintln(markFlag.Output(), "Usage: task-cli-go mark <id> <status>")
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Expected a subcommand")
		os.Exit(2)
	}
	subcommand := os.Args[1]

	switch subcommand {
	case "create":
		if err := createFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if createFlag.NArg() != 1 {
			createFlag.Usage()
			os.Exit(2)
		}
		if err := create(createFlag.Arg(0), filename); err != nil {
			fmt.Fprintf(os.Stderr, "create: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := listFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if listFlag.NArg() != 0 {
			listFlag.Usage()
			os.Exit(2)
		}
		if err := list(filename); err != nil {
			fmt.Fprintf(os.Stderr, "list: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		if err := deleteFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if deleteFlag.NArg() != 1 {
			deleteFlag.Usage()
			os.Exit(2)
		}
		if err := delete(deleteFlag.Arg(0), filename); err != nil {
			fmt.Fprintf(os.Stderr, "delete: %v\n", err)
			os.Exit(1)
		}
	case "update":
		if err := updateFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if updateFlag.NArg() != 2 {
			updateFlag.Usage()
			os.Exit(2)
		}
		if err := update(updateFlag.Arg(0), updateFlag.Arg(1), filename); err != nil {
			fmt.Fprintf(os.Stderr, "update: %v\n", err)
			os.Exit(1)
		}

	case "mark":
		if err := markFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if markFlag.NArg() != 2 {
			markFlag.Usage()
			os.Exit(2)
		}
		status := TaskStatus(markFlag.Arg(1))
		if status != TaskStatusTodo && status != TaskStatusInProgress && status != TaskStatusDone {
			fmt.Fprintf(os.Stderr, "mark: invalid status: %s\n", markFlag.Arg(1))
			fmt.Fprintln(os.Stderr, "Valid statuses: todo, in-progress, done")
			os.Exit(2)
		}
		if err := mark(markFlag.Arg(0), status, filename); err != nil {
			fmt.Fprintf(os.Stderr, "mark: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommand)
		os.Exit(2)
	}
}
