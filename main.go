package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string
	Message string
}

func homehandler(w http.ResponseWriter, r *http.Request) {

	PageData := PageData{
		Title:   "Home Page",
		Message: "Hello World!",
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, PageData)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	http.HandleFunc("/", homehandler)
	fmt.Println("Server is runnning at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
