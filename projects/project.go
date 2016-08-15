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

package projects

import (
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"
)

var prjs map[string]Project = make(map[string]Project)

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

// Project type definition
type Project struct {
	Name         string
	CreationDate time.Time
	Description  string
}

var currentTime = time.Now

// Encodes the project to json
func encode(w io.Writer, p Project) error {
	enc := json.NewEncoder(w)
	return enc.Encode(p)
}

// Saves a Project, returns an error if something has gone wrong.
func Save(p Project) error {
	p.CreationDate = currentTime()
	prjs[p.Name] = p
	return nil
}

// Deletes a project by name, returns an error if the project has not been found.
func Delete(name string) error {
	if _, present := prjs[name]; present {
		delete(prjs, name)
	} else {
		return errors.New("not found: " + name)
	}
	return nil
}

// Returns all the projects sorted by creation date.
func All() []Project {
	ps := make([]Project, len(prjs), len(prjs))
	i := 0
	for _, v := range prjs {
		ps[i] = v
		i++
	}
	sort.Sort(ByCreationDate(ps))
	return ps
}

// Returns a project by name, an error if the project has not been found.
func Get(name string) (error, Project) {
	if v, present := prjs[name]; present {
		return nil, v
	} else {
		return errors.New("not found: " + name), Project{}
	}
}
