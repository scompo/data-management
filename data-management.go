package main

import (
	"flag"
	"github.com/scompo/data-management/domain"
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

func createMockedProjects() []domain.Project {
	prjs := make([]domain.Project, 3)
	tt := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	prjs[0] = domain.Project{
		"test project 1",
		tt,
	}
	prjs[1] = domain.Project{
		"test project 2",
		tt,
	}
	prjs[2] = domain.Project{
		"test project 3",
		tt,
	}
	return prjs
}

var port = flag.String("port", "8080", "server port")

func main() {

	flag.Parse()

	log.Printf("Initializing %v...\n", appName)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/projects", projectsHandler)

	log.Printf("listening on port %v...", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called projectsHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/projects.html")
	p := Page{Title: appName, PageName: "projects"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"Page":     p,
		"Projects": createMockedProjects(),
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called mainHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/index.html")
	p := Page{Title: appName, PageName: "index"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"Page": p,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}