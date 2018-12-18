package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondIfError(t *testing.T) {
	w := httptest.NewRecorder()
	RespondIfError(errors.New("test"), w)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong code returned: %d", w.Code)
	}
}

func TestRespondIfErrorNone(t *testing.T) {
	w := httptest.NewRecorder()
	RespondIfError(nil, w)
	if w.Code != 200 {
		t.Errorf("wrong code returned: %d", w.Code)
	}
}

func TestFailIfError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FailIfError should panic on error")
		}
	}()
	FailIfError(errors.New("test"))
}

func TestFailIfErrorNone(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("FailIfError should not panic without error")
		}
	}()
	FailIfError(nil)
}
