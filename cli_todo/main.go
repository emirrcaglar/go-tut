package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var fileName = "cli_todo/data.json"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// New data structure to match the desired JSON format
type TaskDatabase map[string][]string

func readTasks(filename string) (TaskDatabase, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(TaskDatabase), nil // Return empty map if file doesn't exist
		}
		return nil, err
	}

	var tasks TaskDatabase
	if len(data) > 0 {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func writeTasks(filename string, tasks TaskDatabase) error {
	data, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  To add task: program username \"task description\"")
		fmt.Println("  To list tasks: program task")
		return
	}

	if os.Args[1] == "task" {
		tasks, err := readTasks(fileName)
		check(err)

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		for username, userTasks := range tasks {
			fmt.Printf("User: %s\n", username)
			for i, task := range userTasks {
				fmt.Printf("  %d. %s\n", i+1, task)
			}
		}
		return
	}

	if len(os.Args) >= 3 {
		tasks, err := readTasks(fileName)
		check(err)

		username := os.Args[1]
		taskDescription := os.Args[2]

		// Append task to user's task list
		tasks[username] = append(tasks[username], taskDescription)

		err = writeTasks(fileName, tasks)
		check(err)
		fmt.Printf("Task added successfully for user '%s'\n", username)
	}
}
