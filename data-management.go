package main

import (
	"flag"
	"github.com/scompo/data-management/domain"
	"html/template"
	"log"
	"net/http"
)

const appName = "data-management"

type WebPage struct {
	Title    string
	PageName string
}

var port = flag.String("port", "8080", "server port")

func main() {

	flag.Parse()

	log.Printf("Initializing %v...\n", appName)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/projects/", projectsHandler)
	http.HandleFunc("/projects/new", newProjectHandler)
	http.HandleFunc("/projects/delete", deleteProjectHandler)
	http.HandleFunc("/projects/view", viewProjectHandler)
	http.HandleFunc("/pages/new", pageNewHandler)

	log.Printf("listening on port %v...", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func pageNewHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called pageNewHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/page-new.html")
	p := WebPage{Title: appName, PageName: "New Page"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"WebPage": p,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func viewProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called viewProjectHandler")
	name := r.URL.Query().Get("Name")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/project-view.html")
	p := WebPage{Title: appName, PageName: "View Project"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err, prj := domain.GetProject(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"WebPage":        p,
		"ProjectInfo": prj,
	})
	return
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called DeleteProjectHandler")
	name := r.URL.Query().Get("Name")
	err := domain.DeleteProject(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/projects", http.StatusFound)
	return
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called newProjectHandler")
	switch r.Method {
	case "POST":
		log.Printf("method: POST\n")
		r.ParseForm()
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		pi := domain.ProjectInfo{
			Project: domain.Project{
				Name: name,
			},
			Description: description,
		}
		domain.SaveProject(pi)
		http.Redirect(w, r, "/projects", http.StatusFound)
		return
	case "GET":
		log.Printf("method: GET\n")
		t, err := template.ParseFiles(
			"templates/main.html",
			"templates/header.html",
			"templates/project-new.html")
		p := WebPage{Title: appName, PageName: "New Project"}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, map[string]interface{}{
			"WebPage": p,
		})
		return
	default:
		log.Printf("DEFAULT\n")
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called projectsHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/project-all.html")
	p := WebPage{Title: appName, PageName: "All Projects"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"WebPage":     p,
		"Projects": domain.AllProjects(),
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called mainHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/index.html")
	p := WebPage{Title: appName, PageName: "Main Page"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"WebPage": p,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
