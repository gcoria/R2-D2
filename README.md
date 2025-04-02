# R2-D2 Todo CLI

A simple command-line todo application built with Go and Cobra.

## Features

- **Add tasks**: Add new tasks to your todo list
- **List tasks**: View all your tasks in a tabular format
- **Complete tasks**: Mark tasks as completed
- **Delete tasks**: Remove tasks from your list
- **Interactive REPL mode**: Use the application in an interactive shell

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/R2-D2.git
cd R2-D2

# Build the application
go build -o r2d2
```

## Usage

### CLI Mode

```bash
# Add a new task
./r2d2 add "Buy groceries"

# List all tasks
./r2d2 list

# Mark a task as complete
./r2d2 complete 1

# Delete a task
./r2d2 delete 1

# Get help
./r2d2 help
```

### REPL Mode

Launch the application without any arguments to enter REPL mode:

```bash
./r2d2
```

In REPL mode, you can use the same commands without the "./r2d2" prefix:

```
> add Buy groceries
> list
> complete 1
> delete 1
> help
> exit
```

## Next Steps

### Connect to a NoSQL Database

Currently, the application stores tasks in a CSV file. Future plans include:

### Feat incomning

- Implement DB(mongodb?) integration for task storage
- Implement a `--secret` flag for the add command to create encrypted tasks
- Use AES-256 encryption for sensitive task information
- Require a password to view secret tasks
- Store encryption keys securely

### Additional Planned Features

- Categories and tags for tasks
- Due dates and reminders
- Priority levels
- Recurring tasks
- Export/import functionality
