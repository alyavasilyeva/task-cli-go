package main

import (
	// "flag"
	"encoding/json/v2"
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

func create(description string, filename string) {
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
	byteValue, err = json.Marshal(tasks)
	check(err)

	// write tasks to file
	err = os.WriteFile(filename, byteValue, 0644)
	check(err)
}

func main() {
	const filename = "tasks.json"

	args := os.Args[1:]

	// create task
	create(args[0], filename)
}
