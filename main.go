package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func create(description string, storage JSONFile) error {
	tasks, err := storage.Load()
	if err != nil {
		return err
	}

	tasks, task, err := createTask(tasks, description)
	if errors.Is(err, ErrTaskAlreadyExists) {
		fmt.Println("Task already exists")
		return nil
	}
	if err != nil {
		return err
	}

	if err := storage.Save(tasks); err != nil {
		return err
	}

	fmt.Println("Task created:", task.ID)
	return nil
}

func deleteTaskByID(id string, storage JSONFile) error {
	tasks, err := storage.Load()
	if err != nil {
		return err
	}

	tasks, err = deleteTask(tasks, id)
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Println("Task not found")
		return nil
	}
	if err != nil {
		return err
	}

	if err := storage.Save(tasks); err != nil {
		return err
	}

	fmt.Println("Task deleted:", id)
	return nil
}

func update(id, description string, storage JSONFile) error {
	tasks, err := storage.Load()
	if err != nil {
		return err
	}

	tasks, err = updateTask(tasks, id, description)
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Println("Task not found")
		return nil
	}
	if err != nil {
		return err
	}

	if err := storage.Save(tasks); err != nil {
		return err
	}

	fmt.Println("Task updated:", id, description)
	return nil
}

func mark(id string, status TaskStatus, storage JSONFile) error {
	tasks, err := storage.Load()
	if err != nil {
		return err
	}

	tasks, err = markTask(tasks, id, status)
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Println("Task not found")
		return nil
	}
	if err != nil {
		return err
	}

	if err := storage.Save(tasks); err != nil {
		return err
	}

	fmt.Println("Task marked:", id, status)
	return nil
}

func list(storage JSONFile, status *TaskStatus) error {
	tasks, err := storage.Load()
	if err != nil {
		return err
	}

	tasks = filterTasks(tasks, status)
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	var taskErrors []error
	fmt.Printf("|%-36s|%-20s|%-12s|%-20s|%-20s|\n", "ID", "Description", "Status", "CreatedAt", "UpdatedAt")
	for _, task := range tasks {
		createdAt, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil {
			taskErrors = append(taskErrors, fmt.Errorf("task %s: invalid created_at %q: %w", task.ID, task.CreatedAt, err))
			continue
		}
		updatedAt, err := time.Parse(time.RFC3339, task.UpdatedAt)
		if err != nil {
			taskErrors = append(taskErrors, fmt.Errorf("task %s: invalid updated_at %q: %w", task.ID, task.UpdatedAt, err))
			continue
		}

		fmt.Printf(
			"|%-36s|%-20s|%-12s|%-20s|%-20s|\n",
			task.ID,
			task.Description,
			task.Status,
			createdAt.Format(time.DateTime),
			updatedAt.Format(time.DateTime),
		)
	}

	return errors.Join(taskErrors...)
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, `Usage: task-cli-go <command> [arguments]

Commands:
  create <description>        Create a task
  list [status]               List tasks, optionally filtered by status
  update <id> <description>   Update a task description
  mark <id> <status>          Change a task status
  delete <id>                 Delete a task
  help                        Show this help

Valid statuses: todo, in-progress, done`)
}

func main() {
	storage := JSONFile{Filename: "tasks.json"}

	createFlag := flag.NewFlagSet("create", flag.ExitOnError)
	createFlag.Usage = func() {
		_, _ = fmt.Fprintln(createFlag.Output(), "Usage: task-cli-go create <description>")
	}

	listFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listFlag.Usage = func() {
		_, _ = fmt.Fprintln(listFlag.Output(), "Usage: task-cli-go list [status]")
		_, _ = fmt.Fprintln(listFlag.Output(), "Valid statuses: todo, in-progress, done")
	}

	deleteFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteFlag.Usage = func() {
		_, _ = fmt.Fprintln(deleteFlag.Output(), "Usage: task-cli-go delete <id>")
	}

	updateFlag := flag.NewFlagSet("update", flag.ExitOnError)
	updateFlag.Usage = func() {
		_, _ = fmt.Fprintln(updateFlag.Output(), "Usage: task-cli-go update <id> <description>")
	}

	markFlag := flag.NewFlagSet("mark", flag.ExitOnError)
	markFlag.Usage = func() {
		_, _ = fmt.Fprintln(markFlag.Output(), "Usage: task-cli-go mark <id> <status>")
	}

	if len(os.Args) < 2 {
		printUsage(os.Stderr)
		os.Exit(2)
	}
	subcommand := os.Args[1]

	switch subcommand {
	case "help", "-h", "--help":
		printUsage(os.Stdout)

	case "create":
		if err := createFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if createFlag.NArg() != 1 {
			createFlag.Usage()
			os.Exit(2)
		}
		if err := create(createFlag.Arg(0), storage); err != nil {
			fmt.Fprintf(os.Stderr, "create: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := listFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if listFlag.NArg() > 1 {
			listFlag.Usage()
			os.Exit(2)
		}

		var statusFilter *TaskStatus
		if listFlag.NArg() == 1 {
			status := TaskStatus(listFlag.Arg(0))
			if !status.Valid() {
				fmt.Fprintf(os.Stderr, "list: invalid status: %s\n", status)
				listFlag.Usage()
				os.Exit(2)
			}
			statusFilter = &status
		}

		if err := list(storage, statusFilter); err != nil {
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
		if err := deleteTaskByID(deleteFlag.Arg(0), storage); err != nil {
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
		if err := update(updateFlag.Arg(0), updateFlag.Arg(1), storage); err != nil {
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
		if !status.Valid() {
			fmt.Fprintf(os.Stderr, "mark: invalid status: %s\n", markFlag.Arg(1))
			_, _ = fmt.Fprintln(os.Stderr, "Valid statuses: todo, in-progress, done")
			os.Exit(2)
		}
		if err := mark(markFlag.Arg(0), status, storage); err != nil {
			fmt.Fprintf(os.Stderr, "mark: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommand)
		printUsage(os.Stderr)
		os.Exit(2)
	}
}
