package data_test

import (
	handlers "goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteInvalidID(t *testing.T) {

	mockDataService := &service.MockDataServiceSuccessful{}
	req, err := http.NewRequest("DELETE", "/data/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "invalid") // * Required for routing *

	rr := httptest.NewRecorder()

	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	if rr.Body.String() != `{"error": "Missconfigured ID."}` {
		t.Errorf("handler returned unexpected body: got %v want empty body", rr.Body.String())
	}
}

func TestDeleteError(t *testing.T) {

	mockDataService := &service.MockDataServiceError{}
	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	if strings.TrimSpace(rr.Body.String()) != `Internal Server error` {
		t.Errorf("handler returned unexpected body: got %v want 'Internal Server error'", rr.Body.String())
	}
}

func TestDeleteNotFound(t *testing.T) {
	mockDataService := &service.MockDataServiceNotFound{}
	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	if rr.Body.String() != `{"error": "Resource not found."}` {
		t.Errorf("handler returned unexpected body: got %v want empty body", rr.Body.String())
	}
}

func TestDeleteSuccessful(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}
	req, err := http.NewRequest("DELETE", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	handlers.DeleteHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body: got %v want empty body", rr.Body.String())
	}
}
