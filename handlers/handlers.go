package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "todo-web/database"
	"todo-web/utils"
)

type PageData struct {
	Title  string
	Todos  []db.Todo
	ListID int
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	PageData := PageData{Title: "TODO App"}
	utils.RenderTemplate(w, "home", PageData)
}

func CreateListHandler(app *db.TodoApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			listName := r.FormValue("list_name")
			if listName == "" {
				return
			}
			listID, err := app.CreateList(listName)
			if err != nil {
				http.Error(w, "Error creating list", http.StatusInternalServerError)
				fmt.Printf("Error creating list: %v\n", err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/todo/%d", listID), http.StatusSeeOther)
			return
		}
		PageData := PageData{Title: "Create List"}
		utils.RenderTemplate(w, "create_list", PageData)
	}
}

func TodoListHandler(app *db.TodoApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//urlからlistIDを取得
		listIDStr := r.URL.Path[len("/todo/"):]
		listID, err := strconv.Atoi(listIDStr)
		if err != nil {
			http.Error(w, "Invalid List ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			action := r.FormValue("action")
			switch action {
			case "add":
				HandleAddTodo(app, r, listID)
			case "toggle":
				HandleToggleTodo(app, r, listID)
			case "delete":
				HandleDeleteTodo(app, r, listID)
			default:
				http.Error(w, "Invalid action", http.StatusBadRequest)
				return
			}
		}
		todos, err := app.GetTodos(listID)
		if err != nil {
			http.Error(w, "Error getting todos", http.StatusInternalServerError)
			fmt.Printf("Error getting todos: %v\n", err)
			return
		}

		// リスト名取得
		listName, err := app.GetListName(listID)
		if err != nil {
			http.Error(w, "Error getting list name", http.StatusInternalServerError)
			fmt.Printf("Error getting list name: %v\n", err)
			return
		}

		PageData := PageData{
			Title:  fmt.Sprintf("TODO List - %s", listName),
			Todos:  todos,
			ListID: listID,
		}
		utils.RenderTemplate(w, "todolist", PageData)
	}
}

func HandleAddTodo(app *db.TodoApp, r *http.Request, listID int) {
	description := r.FormValue("description")
	if description == "" {
		return
	}
	err := app.AddTodo(description, listID)
	if err != nil {
		fmt.Printf("Error adding todo: %v\n", err)
	}
}

func HandleToggleTodo(app *db.TodoApp, r *http.Request, listID int) {
	toggleID := r.FormValue("toggle")
	if toggleID == "" {
		return
	}
	id, err := strconv.Atoi(toggleID)
	if err != nil {
		fmt.Printf("Error converting toggleID to int: %v\n", err)
		return
	}
	err = app.ToggleTodo(id, listID)
	if err != nil {
		fmt.Printf("Error toggling todo: %v\n", err)
	}
}

func HandleDeleteTodo(app *db.TodoApp, r *http.Request, listID int) {
	deleteID := r.FormValue("delete")
	if deleteID == "" {
		return
	}
	id, err := strconv.Atoi(deleteID)
	if err != nil {
		fmt.Printf("Error converting deleteID to int: %v\n", err)
		return
	}
	err = app.DeleteTodo(id, listID)
	if err != nil {
		fmt.Printf("Error deleting todo: %v\n", err)
	}
}
