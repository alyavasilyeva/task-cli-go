package main

import (
	"encoding/json/v2"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
	"uuid"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func loadTasks(filename string) []Task {
	// open file or create it if it doesn't exist
	jsonFile, err := os.Open(filename)
	if err != nil && os.IsNotExist(err) {
		jsonFile, err = os.Create(filename)
		check(err)
	}
	defer jsonFile.Close()

	// read file
	byteValue, err := io.ReadAll(jsonFile)
	check(err)
	var tasks []Task

	// unmarshal file into tasks
	json.Unmarshal(byteValue, &tasks)

	return tasks
}

func create(description string, filename string) {
	tasks := loadTasks(filename)

	// check if task already exists
	index := slices.IndexFunc(tasks, func(t Task) bool { return t.Description == description })
	if index != -1 {
		fmt.Println("Task already exists")
		return
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
	check(err)

	// write tasks to file
	err = os.WriteFile(filename, byteValue, 0644)
	check(err)

	fmt.Println("Task created:", task.Id)
}

func list(filename string) {
	tasks := loadTasks(filename)

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	fmt.Printf("|%-36s|%-20s|%-10s|%-20s|%-20s|\n", "ID", "Description", "Status", "CreatedAt", "UpdatedAt")
	for _, task := range tasks {
		createdAt, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing createdAt: %v\n", err)
			continue
		}
		updatedAt, err := time.Parse(time.RFC3339, task.UpdatedAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing updatedAt: %v\n", err)
			continue
		}

		fmt.Printf("|%-36s|%-20s|%-10s|%-20s|%-20s|\n", task.Id, task.Description, task.Status, createdAt.Format(time.DateTime), updatedAt.Format(time.DateTime))
	}
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
		create(createFlag.Arg(0), filename)

	case "list":
		if err := listFlag.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		if listFlag.NArg() != 0 {
			listFlag.Usage()
			os.Exit(2)
		}
		list(filename)

	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommand)
		os.Exit(2)
	}
}
