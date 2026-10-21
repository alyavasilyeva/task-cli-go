# Task CLI

A small Go command-line tool that stores tasks in `tasks.json`.

## Build

Requires Go 1.27 or later.

```sh
go build .
```

## Usage

Show all available commands:

```sh
./task-cli-go --help
```

Create a task:

```sh
./task-cli-go create "Buy milk"
```

List all tasks:

```sh
./task-cli-go list
```

List tasks by status:

```sh
./task-cli-go list todo
./task-cli-go list in-progress
./task-cli-go list done
```

Update a task description:

```sh
./task-cli-go update <id> "Buy oat milk"
```

Mark a task with a new status:

```sh
./task-cli-go mark <id> in-progress
./task-cli-go mark <id> done
```

Delete a task:

```sh
./task-cli-go delete <id>
```

Descriptions containing spaces must be quoted. Valid statuses are `todo`,
`in-progress`, and `done`.
