package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "todo-web/database"
	"todo-web/utils"
)

type PageData struct {
	Title string
	Todos []db.Todo
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	PageData := PageData{Title: "TODO App"}
	utils.RenderTemplate(w, "home", PageData)
}

func TodoListHandler(app *db.TodoApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			HandleAddTodo(app, r)
			HandleToggleTodo(app, r)
			HandleDeleteTodo(app, r)
		}
		todos, err := app.GetTodos()
		if err != nil {
			http.Error(w, "Error getting todos", http.StatusInternalServerError)
			fmt.Printf("Error getting todos: %v\n", err)
			return
		}

		PageData := PageData{Title: "TODO App", Todos: todos}
		utils.RenderTemplate(w, "todolist", PageData)
	}
}

func HandleAddTodo(app *db.TodoApp, r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		return
	}
	err := app.AddTodo(description)
	if err != nil {
		fmt.Printf("Error adding todo: %v\n", err)
	}
}

func HandleToggleTodo(app *db.TodoApp, r *http.Request) {
	toggleID := r.FormValue("toggle")
	if toggleID == "" {
		return
	}
	id, err := strconv.Atoi(toggleID)
	if err != nil {
		fmt.Printf("Error converting toggleID to int: %v\n", err)
		return
	}
	err = app.ToggleTodo(id)
	if err != nil {
		fmt.Printf("Error toggling todo: %v\n", err)
	}
}

func HandleDeleteTodo(app *db.TodoApp, r *http.Request) {
	deleteID := r.FormValue("delete")
	if deleteID == "" {
		return
	}
	id, err := strconv.Atoi(deleteID)
	if err != nil {
		fmt.Printf("Error converting deleteID to int: %v\n", err)
		return
	}
	err = app.DeleteTodo(id)
	if err != nil {
		fmt.Printf("Error deleting todo: %v\n", err)
	}
}
