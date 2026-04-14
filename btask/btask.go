package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

const (
	DataDir  = "/var/lib/btask"
	DataFile = "/var/lib/btask/tasks.json"
	Cyan     = "\033[1;36m"
	Green    = "\033[1;32m"
	Yellow   = "\033[1;33m"
	Gray     = "\033[0;90m"
	Reset    = "\033[0m"
	Bold     = "\033[1m"
)

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(DataFile); os.IsNotExist(err) {
		return []Task{}, nil
	}
	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	_ = os.MkdirAll(DataDir, 0755)
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(DataFile, data, 0644)
}

func main() {
	tasks, _ := loadTasks()
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("\n  %s%s[ BLOCK TASKS ]%s\n", Cyan, Bold, Reset)
		fmt.Println("  ---------------------------")
		if len(tasks) == 0 {
			fmt.Println("  (No pending tasks. Relax!)")
		}
		for _, t := range tasks {
			status := "[ ]"
			color := Yellow
			text := t.Text
			if t.Done {
				status = "[X]"
				color = Gray
				text = Gray + t.Text + Reset
			}
			fmt.Printf("  %s%d. %s %s%s\n", Cyan, t.ID, status, color, text, Reset)
		}
		fmt.Println("\n  Usage: btask add \"msg\" | done <id> | clear")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: Missing task description")
			return
		}
		id := 1
		if len(tasks) > 0 {
			id = tasks[len(tasks)-1].ID + 1
		}
		tasks = append(tasks, Task{ID: id, Text: args[1], Done: false})
		saveTasks(tasks)
		fmt.Printf("  %s✔ Added task: %s%s\n", Green, args[1], Reset)

	case "done":
		if len(args) < 2 {
			fmt.Println("Error: Missing task ID")
			return
		}
		id, _ := strconv.Atoi(args[1])
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				found = true
				break
			}
		}
		if found {
			saveTasks(tasks)
			fmt.Printf("  %s✔ Task %d marked as done!%s\n", Green, id, Reset)
		} else {
			fmt.Println("Error: Task ID not found")
		}

	case "clear":
		var active []Task
		for _, t := range tasks {
			if !t.Done {
				active = append(active, t)
			}
		}
		saveTasks(active)
		fmt.Printf("  %s✔ Cleared completed tasks.%s\n", Green, Reset)

	case "reset":
		saveTasks([]Task{})
		fmt.Println("  ✔ Task list reset.")

	default:
		fmt.Println("Unknown command. Try: add, done, clear, reset")
	}
}
