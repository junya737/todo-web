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
	createTodosTableQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL,
		list_id INTEGER NOT NULL,
		FOREIGN KEY (list_id) REFERENCES lists(id)
	);`
	_, err := app.DB.Exec(createTodosTableQuery)
	if err != nil {
		return err
	}

	// リストテーブルの作成
	createListsTableQuery := `
	CREATE TABLE IF NOT EXISTS lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	_, err = app.DB.Exec(createListsTableQuery)
	return err
}

func (app *TodoApp) GetTodos(listID int) ([]Todo, error) {
	getTodoQuery := "SELECT id, description, completed FROM todos WHERE list_id = ?"
	rows, err := app.DB.Query(getTodoQuery, listID)
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

func (app *TodoApp) AddTodo(description string, listID int) error {
	insertTodoQuery := "INSERT INTO todos (description, completed, list_id) VALUES (?, ?, ?)"
	_, err := app.DB.Exec(insertTodoQuery, description, false, listID)
	return err
}

func (app *TodoApp) ToggleTodo(id int, listID int) error {
	toggleTodoQuery := "UPDATE todos SET completed = NOT completed WHERE id = ? AND list_id = ?"
	_, err := app.DB.Exec(toggleTodoQuery, id, listID)
	return err
}

func (app *TodoApp) DeleteTodo(id int, listID int) error {
	deleteTodoQuery := "DELETE FROM todos WHERE id = ? AND list_id = ?"
	_, err := app.DB.Exec(deleteTodoQuery, id, listID)
	return err
}
