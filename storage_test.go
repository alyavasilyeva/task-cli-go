package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestJSONFileLoadMissingFile(t *testing.T) {
	storage := JSONFile{Filename: filepath.Join(t.TempDir(), "tasks.json")}

	tasks, err := storage.Load()

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if tasks != nil {
		t.Errorf("Load() returned %#v, want nil", tasks)
	}
}

func TestJSONFileLoadEmptyFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(filename, nil, taskFileMode); err != nil {
		t.Fatalf("write empty file: %v", err)
	}
	storage := JSONFile{Filename: filename}

	tasks, err := storage.Load()

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if tasks != nil {
		t.Errorf("Load() returned %#v, want nil", tasks)
	}
}

func TestJSONFileLoadInvalidJSON(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(filename, []byte("{invalid"), taskFileMode); err != nil {
		t.Fatalf("write invalid JSON: %v", err)
	}
	storage := JSONFile{Filename: filename}

	if _, err := storage.Load(); err == nil {
		t.Fatal("Load() error = nil, want decoding error")
	}
}

func TestJSONFileSaveAndLoad(t *testing.T) {
	directory := t.TempDir()
	filename := filepath.Join(directory, "tasks.json")
	storage := JSONFile{Filename: filename}
	want := []Task{{
		ID:          "0199f00d-0000-7000-8000-000000000001",
		Description: "Buy milk",
		Status:      TaskStatusTodo,
		CreatedAt:   "2026-10-21T10:00:00Z",
		UpdatedAt:   "2026-10-21T10:00:00Z",
	}}

	if err := storage.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := storage.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() returned %#v, want %#v", got, want)
	}

	matches, err := filepath.Glob(filepath.Join(directory, ".tasks.json.tmp-*"))
	if err != nil {
		t.Fatalf("find temporary files: %v", err)
	}
	if len(matches) != 0 {
		t.Errorf("Save() left temporary files: %v", matches)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("stat saved file: %v", err)
	}
	if gotMode := info.Mode().Perm(); gotMode != taskFileMode {
		t.Errorf("Save() file mode = %o, want %o", gotMode, taskFileMode)
	}
}

func TestJSONFileSaveCleansUpAfterRenameFailure(t *testing.T) {
	directory := t.TempDir()
	filename := filepath.Join(directory, "tasks.json")
	if err := os.Mkdir(filename, 0o755); err != nil {
		t.Fatalf("create destination directory: %v", err)
	}
	storage := JSONFile{Filename: filename}

	if err := storage.Save([]Task{{ID: "1"}}); err == nil {
		t.Fatal("Save() error = nil, want rename error")
	}

	matches, err := filepath.Glob(filepath.Join(directory, ".tasks.json.tmp-*"))
	if err != nil {
		t.Fatalf("find temporary files: %v", err)
	}
	if len(matches) != 0 {
		t.Errorf("Save() left temporary files after failure: %v", matches)
	}
}
