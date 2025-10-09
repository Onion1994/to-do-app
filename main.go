package main

import (
	"flag"
)

func main() {
	todos, _ := LoadTodos()

	viewFlag := flag.Bool("view", false, "View to-do list")
	addFlag := flag.String("add", "", "Add a new to-do item")
	findFlag := flag.String("find", "", "Find to-do item by description")
	updateStatusFlag := flag.String("update-status", "", "Update a to-do item status")
	updateDescriptionFlag := flag.String("update-description", "", "Update a to-do item description")
	removeFlag := flag.String("remove", "", "Remove a to-do item")

	flag.Parse()

	if *viewFlag {
		PrintTodos(todos)
	}

	if *addFlag != "" {
		todos = AddNewItem(todos, *addFlag)
		SaveTodos(todos)
	}

	if *removeFlag != "" {
		todos = RemoveItem(todos, *removeFlag)
		SaveTodos(todos)
	}

	if *findFlag != "" && *updateStatusFlag != "" {
		UpdateStatus(todos, *findFlag, *updateStatusFlag)
		SaveTodos(todos)
	}

	if *findFlag != "" && *updateDescriptionFlag != "" {
		UpdateDesc(todos, *findFlag, *updateDescriptionFlag)
		SaveTodos(todos)
	}
}
