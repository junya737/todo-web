package main

import (
	"database/sql"
	"fmt"
	"net/http"

	db "todo-web/database"
	"todo-web/handlers"
)

func main() {
	// データベース接続
	dbConn, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer dbConn.Close()

	app := &db.TodoApp{DB: dbConn}

	err = app.InitDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/todo/", handlers.TodoListHandler(app))
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
