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

type Task struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"` // todo, in-progress, done
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
		Status:      "todo",
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

func main() {
	const filename = "tasks.json"

	createFlag := flag.NewFlagSet("create", flag.ExitOnError)
	createFlag.Usage = func() {
		fmt.Fprintln(createFlag.Output(), "Usage: task-cli-go create <description>")
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
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommand)
		os.Exit(2)
	}
}
