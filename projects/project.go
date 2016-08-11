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

var prjs map[string]ProjectInfo = make(map[string]ProjectInfo)

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

type Project struct {
	Name         string
	CreationDate time.Time
}

var currentTime = time.Now

func (p *Project) Encode(w io.Writer) {
	p.CreationDate = currentTime()
	enc := json.NewEncoder(w)
	enc.Encode(p)
	return
}

// Saves a ProjectInfo
func Save(p ProjectInfo) error {
	p.CreationDate = currentTime()
	prjs[p.Name] = p
	return nil
}

// Deletes a project by name
func Delete(name string) error {
	if _, present := prjs[name]; present {
		delete(prjs, name)
	} else {
		return errors.New("not found: " + name)
	}
	return nil
}

func All() []ProjectInfo {
	ps := make([]ProjectInfo, len(prjs), len(prjs))
	i := 0
	for _, v := range prjs {
		ps[i] = v
		i++
	}
	sort.Sort(ByCreationDate(ps))
	return ps
}

func Get(name string) (error, ProjectInfo) {
	if v, present := prjs[name]; present {
		return nil, v
	} else {
		return errors.New("not found: " + name), ProjectInfo{}
	}
}
