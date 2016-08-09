package domain

import (
	"encoding/json"
	"github.com/minio/minio-go"
	"io"
	"time"
)

var projects []Project = make([]Project, 0)

type ProjectInfo struct {
	Project
	Description string
}

type Project minio.BucketInfo

var currentTime = time.Now

func (p *Project) Encode(w io.Writer) {
	p.CreationDate = currentTime()
	enc := json.NewEncoder(w)
	enc.Encode(p)
	return
}

func SaveProject(p ProjectInfo) error {
	p.CreationDate = currentTime()
	projects = append(projects, p.Project)
	return nil
}

func AllProjects() []Project {
	return projects
}
