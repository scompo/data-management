package utils

import (
	"errors"
	"net/http"
	"testing"
)

func failFunction(w http.ResponseWriter, r *http.Request) error {
	return errors.New("I always do fail!")
}

func TestFailFunction(t *testing.T) {
	err := failFunction(nil, nil)
	if err == nil {
		t.Errorf("should return an error")
	}
}

func TestServeHttp(t *testing.T) {
	t.Errorf("todo")
}
