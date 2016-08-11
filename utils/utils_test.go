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
