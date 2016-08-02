package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
)

const appName = "data-management"

type Page struct {
	Title    string
	PageName string
}

type Project struct {
	Name         string
	Description  string
	CreationDate time.Time
	LastEdit     time.Time
}

var port = flag.String("port", "8080", "server port")

func main() {

	flag.Parse()

	log.Printf("Initializing %v...\n", appName)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mainHandler)

	log.Printf("listening on port %v...", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called mainHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/index.html")
	p := Page{Title: appName, PageName: "index"}
	prjs := make([]Project, 1)
	tt := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	prjs[0] = Project{
		"test project",
		"test project description",
		tt,
		tt,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"Page":     p,
		"Projects": prjs,
	})
}
