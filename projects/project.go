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

// Package projects contains projects definition and functions.
package projects

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PrjDir is the directory where the projects are saved in.
var PrjDir string

var prjIndexName = "projects.json"

type byCreationDate []Project

func (p byCreationDate) Len() int {
	return len(p)
}

func (p byCreationDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p byCreationDate) Less(i, j int) bool {
	return p[i].CreationDate.Before(p[j].CreationDate)
}

// Project type definition
type Project struct {
	Name         string
	CreationDate time.Time
	Description  string
}

var currentTime = time.Now

// Save saves a Project.
// Returns an error if something has gone wrong.
func Save(p Project) error {
	if Exists(p.Name) {
		return errors.New("project name already existent: " + p.Name)
	}
	err := createProjectDir(p.Name)
	if err != nil {
		return err
	}
	p.CreationDate = currentTime()
	persist(p)
	return nil
}

func persist(p Project) error {
	projects, err := deserialize()
	if err != nil {
		return err
	}
	projects = append(projects, p)
	return serialize(projects)
}

func deserialize() ([]Project, error) {
	r, err := os.Open(filepath.Join(PrjDir, prjIndexName))
	var data []Project
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	dec := json.NewDecoder(r)
	err = dec.Decode(&data)
	return data, err
}

func serialize(prjs []Project) error {
	w, err := os.Create(filepath.Join(PrjDir, prjIndexName))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	return enc.Encode(prjs)
}

// GetProjectPath returns the base path for a project.
func GetProjectPath(name string) string {
	return filepath.Join(PrjDir, name)
}

func createProjectDir(name string) error {
	return os.MkdirAll(GetProjectPath(name), 0775)
}

func deleteProjectDir(name string) error {
	return os.RemoveAll(GetProjectPath(name))
}

// Delete deletes a project by name.
func Delete(name string) error {
	if Exists(name) {
		ps, err := deserialize()
		if err != nil {
			return err
		}
		ind := -1
		for i, v := range ps {
			if v.Name == name {
				ind = i
				break
			}
		}
		ps = append(ps[:ind], ps[ind+1:]...)
		err = serialize(ps)
		if err != nil {
			return err
		}
		return deleteProjectDir(name)
	}
	return nil
}

// All returns all the projects sorted by creation date.
func All() []Project {
	ps, err := deserialize()
	if err != nil {
		ps = make([]Project, 0)
	} else {
		sort.Sort(byCreationDate(ps))
	}
	return ps
}

// Get returns a project by name.
func Get(name string) (Project, error) {
	if !Exists(name) {
		return Project{}, errors.New("not present")
	}
	prjs := All()
	for _, prj := range prjs {
		if prj.Name == name {
			return prj, nil
		}
	}
	return Project{}, errors.New("should not get here")
}

// Exists checks if a project exists.
func Exists(name string) bool {
	prjs := All()
	for _, prj := range prjs {
		if prj.Name == name {
			return true
		}
	}
	return false
}
