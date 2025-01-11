package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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

var todos = []Todo{
	{ID: 1, Description: "Learn Go", Completed: false},
	{ID: 2, Description: "Build a TODO app", Completed: true},
}

var nextID = 3

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleAddTodo(r)
		handleToggleTodo(r)
		handleDeleteTodo(r)
	}

	PageData := PageData{Title: "TODO App", Todos: todos}
	RenderTemplate(w, "home", PageData)
}

func handleAddTodo(r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		return
	}
	newTodo := Todo{
		ID:          nextID,
		Description: description,
		Completed:   false,
	}
	nextID++
	todos = append(todos, newTodo)
}

func handleToggleTodo(r *http.Request) {
	toggleID := r.FormValue("toggle")
	if toggleID == "" {
		return
	}
	id, err := strconv.Atoi(toggleID)
	if err != nil {
		fmt.Printf("Error converting toggleID to int: %v\n", err)
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = !todo.Completed
			break
		}
	}

}

func handleDeleteTodo(r *http.Request) {
	deleteID := r.FormValue("delete")

	if deleteID == "" {
		return
	}
	id, err := strconv.Atoi(deleteID)
	if err != nil {
		fmt.Printf("Error converting deleteID to int: %v\n", err)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
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
	http.HandleFunc("/", homeHandler)
	fmt.Println("Server is runnning at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
