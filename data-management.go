package main

import (
	"errors"
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

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	flag.Parse()

	log.Printf("Initializing %v...\n", appName)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", appHandler(mainHandler))
	http.Handle("/projects/", appHandler(projectsHandler))
	http.Handle("/projects/new", appHandler(newProjectHandler))
	http.Handle("/projects/delete", appHandler(deleteProjectHandler))
	http.Handle("/projects/view", appHandler(viewProjectHandler))
	http.Handle("/pages/new", appHandler(pageNewHandler))

	log.Printf("listening on port %v...", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func pageNewHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called viewProjectHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/pages/new.html")
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "New Page",
		},
	})
}

func viewProjectHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called viewProjectHandler")
	name := r.URL.Query().Get("Name")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/projects/view.html")
	if err != nil {
		return err
	}
	err, prj := domain.GetProject(name)
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "View Project",
		},
		"ProjectInfo": prj,
	})
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called DeleteProjectHandler")
	name := r.URL.Query().Get("Name")
	err := domain.DeleteProject(name)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/projects", http.StatusFound)
	return nil
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called newProjectHandler")
	switch r.Method {
	case "POST":
		log.Printf("method: POST\n")
		err := r.ParseForm()
		if err != nil {
			return err
		}
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
		return nil
	case "GET":
		log.Printf("method: GET\n")
		t, err := template.ParseFiles(
			"templates/main.html",
			"templates/header.html",
			"templates/projects/new.html")
		if err != nil {
			return err
		}
		return t.Execute(w, map[string]interface{}{
			"WebPage": WebPage{
				Title:    appName,
				PageName: "New Project",
			},
		})
	default:
		log.Printf("DEFAULT\n")
		return errors.New("method not supported, " + r.Method)
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called projectsHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/projects/list.html")
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "All Projects",
		},
		"Projects": domain.AllProjects(),
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Called mainHandler")
	t, err := template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		"templates/index.html")
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "Main Page",
		},
	})
}
