package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int
	Description string
	Completed   bool
}

type PageData struct {
	Title string
	Todos []Todo
}

type TodoApp struct {
	DB *sql.DB
}

func (app *TodoApp) initDB() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`
	_, err := app.DB.Exec(createTableQuery)
	return err
}

func (app *TodoApp) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.handleAddTodo(r)
		app.handleToggleTodo(r)
		app.handleDeleteTodo(r)
	}
	todos, err := app.getTodos()
	if err != nil {
		http.Error(w, "Error getting todos", http.StatusInternalServerError)
		fmt.Printf("Error getting todos: %v\n", err)
		return
	}

	PageData := PageData{Title: "TODO App", Todos: todos}
	RenderTemplate(w, "home", PageData)
}

func (app *TodoApp) getTodos() ([]Todo, error) {
	rows, err := app.DB.Query("SELECT id, description, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (app *TodoApp) addTodo(description string) error {
	insertTodoQuery := "INSERT INTO todos (description, completed) VALUES (?, ?)"
	_, err := app.DB.Exec(insertTodoQuery, description, false)
	return err
}

func (app *TodoApp) toggleTodo(id int) error {
	toggleTodoQuery := "UPDATE todos SET completed = NOT completed WHERE id = ?"
	_, err := app.DB.Exec(toggleTodoQuery, id)
	return err
}

func (app *TodoApp) deleteTodo(id int) error {
	deleteTodoQuery := "DELETE FROM todos WHERE id = ?"
	_, err := app.DB.Exec(deleteTodoQuery, id)
	return err
}

func (app *TodoApp) handleAddTodo(r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		return
	}
	err := app.addTodo(description)
	if err != nil {
		fmt.Printf("Error adding todo to database: %v\n", err)
	}
}

func (app *TodoApp) handleToggleTodo(r *http.Request) {
	toggleID := r.FormValue("toggle")
	if toggleID == "" {
		return
	}
	id, err := strconv.Atoi(toggleID)
	if err != nil {
		fmt.Printf("Error converting toggleID to int: %v\n", err)
		return
	}
	err = app.toggleTodo(id)
	if err != nil {
		fmt.Printf("Error toggling todo in database: %v\n", err)
	}
}

func (app *TodoApp) handleDeleteTodo(r *http.Request) {
	deleteID := r.FormValue("delete")

	if deleteID == "" {
		return
	}
	id, err := strconv.Atoi(deleteID)
	if err != nil {
		fmt.Printf("Error converting deleteID to int: %v\n", err)
		return
	}
	err = app.deleteTodo(id)
	if err != nil {
		fmt.Printf("Error deleting todo from database: %v\n", err)
	}

}

func RenderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func main() {

	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()
	app := TodoApp{DB: db}

	err = app.initDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	http.HandleFunc("/", app.homeHandler)
	fmt.Println("Server is runnning at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
