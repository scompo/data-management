package domain

import (
	"encoding/json"
	"errors"
	"github.com/minio/minio-go"
	"io"
	"sort"
	"time"
)

var projects []Project = make([]Project, 0)

type ProjectInfo struct {
	Project
	Description string
}

type ByCreationDate []Project

func (p ByCreationDate) Len() int {
	return len(p)
}

func (p ByCreationDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByCreationDate) Less(i, j int) bool {
	return p[i].CreationDate.Before(p[j].CreationDate)
}

type Project minio.BucketInfo

var currentTime = time.Now

func (p *Project) Encode(w io.Writer) {
	p.CreationDate = currentTime()
	enc := json.NewEncoder(w)
	enc.Encode(p)
	return
}

// Saves a ProjectInfo
func SaveProject(p ProjectInfo) error {
	p.CreationDate = currentTime()
	projects = append(projects, p.Project)
	sort.Sort(ByCreationDate(projects))
	return nil
}

// Deletes a project by name
func DeleteProject(name string) error {
	for i, p := range projects {
		if p.Name == name {
			projects = append(projects[:i], projects[i+1:]...)
			return nil
		}
	}
	return errors.New("not found: " + name)
}

func AllProjects() []Project {
	return projects
}
