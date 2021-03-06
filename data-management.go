/*
Copyright (c) 2016, Mauro Scomparin
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of data-management nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"errors"
	"flag"
	"github.com/scompo/data-management/projects"
	"github.com/scompo/data-management/utils"
	"html/template"
	"log"
	"net/http"
	"os"
)

const appName = "data-management"

type WebPage struct {
	Title    string
	PageName string
}

func main() {

	conf := utils.CreateConfig("port", "prj-dir")

	conf["port"] = flag.String("port", "8080", "server port")
	conf["prj-dir"] = flag.String("prj-dir", "data/projects", "project directory path")

	flag.Parse()

	err := initialize(conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func initialize(conf utils.Config) error {

	log.Printf("Initializing...\n")

	for k, v := range conf {
		log.Printf("%v: %v\n", k, *v)
	}

	projects.PrjDir = *conf["prj-dir"]
	
	err := os.MkdirAll(projects.PrjDir, 0775)
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", utils.AppHandler(mainHandler))
	http.Handle("/projects/", utils.AppHandler(projectsHandler))
	http.Handle("/projects/new", utils.AppHandler(newProjectHandler))
	http.Handle("/projects/delete", utils.AppHandler(deleteProjectHandler))
	http.Handle("/projects/view", utils.AppHandler(viewProjectHandler))
	http.Handle("/pages/new", utils.AppHandler(pageNewHandler))

	err = http.ListenAndServe(":"+*conf["port"], nil)
	if err != nil {
		return err
	}

	log.Printf("done\n")
	return nil
}

func pageNewHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := prepareAppTemplate("templates/pages/new.html")
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
	name := r.URL.Query().Get("Name")
	t, err := prepareAppTemplate("templates/projects/view.html")
	if err != nil {
		return err
	}
	prj, err := projects.Get(name)
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "View Project",
		},
		"Project": prj,
	})
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("Name")
	err := projects.Delete(name)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/projects", http.StatusFound)
	return nil
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			return err
		}
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		err = projects.Save(projects.Project{
			Name:        name,
			Description: description,
		})
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/projects", http.StatusFound)
		return nil
	case "GET":
		t, err := prepareAppTemplate("templates/projects/new.html")
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
		return errors.New("method not supported, " + r.Method)
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := prepareAppTemplate("templates/projects/list.html")
	if err != nil {
		return err
	}
	return t.Execute(w, map[string]interface{}{
		"WebPage": WebPage{
			Title:    appName,
			PageName: "All Projects",
		},
		"Projects": projects.All(),
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := prepareAppTemplate("templates/index.html")
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

func prepareAppTemplate(contentTemplate string) (*template.Template, error) {
	return template.ParseFiles(
		"templates/main.html",
		"templates/header.html",
		contentTemplate)
}
