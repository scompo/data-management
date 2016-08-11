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

package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func failFunction(w http.ResponseWriter, r *http.Request) error {
	return errors.New("fail")
}

func TestFailFunction(t *testing.T) {
	err := failFunction(nil, nil)
	if err == nil {
		t.Errorf("should return an error\n")
	}
}

func TestServeHTTP(t *testing.T) {
	w := httptest.NewRecorder()
	AppHandler(failFunction).ServeHTTP(w, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("should return %v but returned %v\n", http.StatusInternalServerError, w.Code)
	}
	if b := w.Body.String(); !strings.Contains(b, "fail") {
		t.Errorf("returned something else in the body: \"%v\"\n", b)
	}
}
