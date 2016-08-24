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
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func setup(t *testing.T) {
	projectDirectory, err := ioutil.TempDir("", "projects")
	if err != nil {
		t.Errorf("error setting test directory")
	}
	PrjDir = projectDirectory
	//t.Logf("set test directory: %v", PrjDir)
	currentTime = func() time.Time {
		return testTime
	}
}

func teardown(t *testing.T) {
	//t.Logf("deleting test directory: %v", PrjDir)
	err := os.RemoveAll(PrjDir)
	if err != nil {
		t.Errorf("error deleting test directory")
	}
	currentTime = time.Now
}

var testTime = time.Now()

func TestExists(t *testing.T) {

	setup(t)

	p := Project{
		Name:        "testName",
		Description: "test description",
	}
	Save(p)
	res := Exists(p.Name)
	if !res {
		t.Errorf("should exist!")
	}
	res = Exists("not existent")
	if res {
		t.Errorf("should not exist!")
	}

	teardown(t)

}

func TestDelete(t *testing.T) {

	setup(t)

	p := Project{
		Name:        "testName",
		Description: "test description",
	}
	err := Save(p)
	if err != nil {
		t.Errorf("Error saving: %v\n", err)
	}
	res := Exists(p.Name)
	if !res {
		t.Errorf("not saved!")
	}
	err = Delete(p.Name)
	if err != nil {
		t.Errorf("Error deleting: %v\n", err)
	}
	res = Exists(p.Name)
	if res {
		t.Errorf("not deleted!")
	}
	res = Exists(p.Name)
	if res {
		t.Errorf("should not error if project not existent!")
	}
	teardown(t)
}

func TestGet(t *testing.T) {

	setup(t)

	p := Project{
		Name:        "testName",
		Description: "test description",
	}
	err := Save(p)
	if err != nil {
		t.Errorf("Error saving: %v\n", err)
	}
	pSaved, err := Get(p.Name)
	if err != nil {
		t.Errorf("Error getting: %v\n", err)
	}
	if p.Name != pSaved.Name {
		t.Errorf("Expected name \"%v\" but was \"%v\"", p.Name, pSaved.Name)
	}
	if p.Description != pSaved.Description {
		t.Errorf("Expected description \"%v\", but was \"%v\"", p.Description, pSaved.Description)
	}

	teardown(t)
}

func TestSave(t *testing.T) {

	setup(t)

	p := Project{
		Name:        "testName",
		Description: "test description",
	}
	err := Save(p)
	if err != nil {
		t.Errorf("Error saving: %v\n", err)
	}
	pSaved, err := Get(p.Name)
	if err != nil {
		t.Errorf("Error getting: %v\n", err)
	}
	if p.Name != pSaved.Name {
		t.Errorf("Expected name \"%v\" but was \"%v\"", p.Name, pSaved.Name)
	}
	if p.Description != pSaved.Description {
		t.Errorf("Expected description \"%v\", but was \"%v\"", p.Description, pSaved.Description)
	}
	if !testTime.Equal(pSaved.CreationDate) {
		t.Errorf("Date should be updated to \"%v\", but was \"%v\"", testTime, pSaved.CreationDate)
	}
	err = Save(p)
	if err == nil {
		t.Errorf("no error for project name already existent\n")
	}
	teardown(t)
}

func TestAll(t *testing.T) {

	setup(t)

	res := All()
	if len(res) != 0 {
		t.Errorf("Nothing should be saved, but found %v projects", len(res))
	}
	p := Project{
		Name:        "testName",
		Description: "test description",
	}
	p2 := Project{
		Name:        "testName2",
		Description: "test description2",
	}
	err := Save(p)
	if err != nil {
		t.Errorf("Error saving: %v\n", err)
	}
	err = Save(p2)
	if err != nil {
		t.Errorf("Error saving: %v\n", err)
	}
	res = All()
	if len(res) != 2 {
		t.Errorf("Saved 2 projects, but found %v", len(res))
	} else {
		if p.Name != res[0].Name || p2.Name != res[1].Name {
			t.Errorf("not in order")
		}
	}

	teardown(t)
}
