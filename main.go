package main

import (
	"fmt"
	"html/template"
	"net/http"
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

func homehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		description := r.FormValue("description")
		if description != "" {
			newTodo := Todo{ID: nextID, Description: description, Completed: false}
			nextID++
			todos = append(todos, newTodo)
		}
	}
	PageData := PageData{Title: "TODO App", Todos: todos}
	RenderTemplate(w, "home", PageData)
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
	}
}

func main() {
	http.HandleFunc("/", homehandler)
	fmt.Println("Server is runnning at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
