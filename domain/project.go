package domain

import "time"

type Project struct {
	Name         string
	Description  string
	CreationDate time.Time
	LastEdit     time.Time
}

type ProjectsRepository struct {
	projects []Project
}

func (pr ProjectsRepository) All() []Project {
	return pr.projects
}

func (pr *ProjectsRepository) Add(prjs []Project) {
	pr.projects = prjs
}
