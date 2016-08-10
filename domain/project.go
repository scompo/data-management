package domain

import (
	"encoding/json"
	"errors"
	"github.com/minio/minio-go"
	"io"
	"sort"
	"time"
)

var projects map[string]ProjectInfo = make(map[string]ProjectInfo)

type ProjectInfo struct {
	Project
	Description string
}

type ByCreationDate []ProjectInfo

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
	projects[p.Name] = p
	return nil
}

// Deletes a project by name
func DeleteProject(name string) error {
	if _, present := projects[name]; present {
		delete(projects, name)
	} else {
		return errors.New("not found: " + name)
	}
	return nil
}

func AllProjects() []ProjectInfo {
	ps := make([]ProjectInfo, len(projects), len(projects))
	i := 0
	for _, v := range projects {
		ps[i] = v
		i++
	}
	sort.Sort(ByCreationDate(ps))
	return ps
}

func GetProject(name string) (error, ProjectInfo) {
	if v, present := projects[name]; present {
		return nil, v
	} else {
		return errors.New("not found: " + name), ProjectInfo{}
	}
}
