package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int
	Description string
	Completed   bool
}

type TodoApp struct {
	DB *sql.DB
}

func (app *TodoApp) InitDB() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`
	_, err := app.DB.Exec(createTableQuery)
	return err
}

func (app *TodoApp) GetTodos() ([]Todo, error) {
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

func (app *TodoApp) AddTodo(description string) error {
	insertTodoQuery := "INSERT INTO todos (description, completed) VALUES (?, ?)"
	_, err := app.DB.Exec(insertTodoQuery, description, false)
	return err
}

func (app *TodoApp) ToggleTodo(id int) error {
	toggleTodoQuery := "UPDATE todos SET completed = NOT completed WHERE id = ?"
	_, err := app.DB.Exec(toggleTodoQuery, id)
	return err
}

func (app *TodoApp) DeleteTodo(id int) error {
	deleteTodoQuery := "DELETE FROM todos WHERE id = ?"
	_, err := app.DB.Exec(deleteTodoQuery, id)
	return err
}
