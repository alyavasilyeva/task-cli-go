package main

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const taskFileMode = 0o644

type JSONFile struct {
	Filename string
}

func (s JSONFile) Load() ([]Task, error) {
	data, err := os.ReadFile(s.Filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", s.Filename, err)
	}
	if len(data) == 0 {
		return nil, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("decode %s: %w", s.Filename, err)
	}

	return tasks, nil
}

func (s JSONFile) Save(tasks []Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	directory := filepath.Dir(s.Filename)
	temporary, err := os.CreateTemp(directory, "."+filepath.Base(s.Filename)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", s.Filename, err)
	}

	temporaryName := temporary.Name()
	closed := false
	defer func() {
		if !closed {
			_ = temporary.Close()
		}
		_ = os.Remove(temporaryName)
	}()

	if err := temporary.Chmod(taskFileMode); err != nil {
		return fmt.Errorf("set permissions on temporary file %s: %w", temporaryName, err)
	}

	written, err := temporary.Write(data)
	if err != nil {
		return fmt.Errorf("write temporary file %s: %w", temporaryName, err)
	}
	if written != len(data) {
		return fmt.Errorf("write temporary file %s: %w", temporaryName, io.ErrShortWrite)
	}

	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary file %s: %w", temporaryName, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file %s: %w", temporaryName, err)
	}
	closed = true

	if err := os.Rename(temporaryName, s.Filename); err != nil {
		return fmt.Errorf("replace %s: %w", s.Filename, err)
	}

	return nil
}
