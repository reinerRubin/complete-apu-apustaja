package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reinerRubin/complete-apu-apustaja/internal/completer"
)

func TestCompleteHandlerWithoutTypes(t *testing.T) {
	request, err := http.NewRequest("GET", "/complete", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler := NewCompleterHandler(completer.NewDummyCompleter())
	handler.Handle(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
}

func TestCompleteHandlerBase(t *testing.T) {
	request, err := http.NewRequest("GET", "/complete?term=vechnayavesna&types[]=city", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler := NewCompleterHandler(completer.NewDummyCompleter())
	handler.Handle(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
}
