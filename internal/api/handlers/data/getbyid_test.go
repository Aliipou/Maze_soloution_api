package data_test

import (
	"encoding/json"
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetByIDInvalidID(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}
	req, err := http.NewRequest("GET", "/data/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "invalid") // * Required for routing *
	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if strings.TrimSpace(rr.Body.String()) != `{"error": "Missconfigured ID."}` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `{"error": "Missconfigured ID."}`)
	}
}

func TestGetByIdInternalError(t *testing.T) {
	mockDataService := &service.MockDataServiceError{}
	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	if strings.TrimSpace(rr.Body.String()) != `Internal server error.` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `Internal server error.`)
	}
}

func TestGetByIdNotFound(t *testing.T) {
	mockDataService := &service.MockDataServiceNotFound{}
	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	if strings.TrimSpace(rr.Body.String()) != `{"error": "Resource not found."}` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `{"error": "Resource not found."}`)
	}
}

func TestGetByIdSuccessful(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}
	req, err := http.NewRequest("GET", "/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "1") // * Required for routing *

	rr := httptest.NewRecorder()

	data.GetByIDHandler(rr, req, log.Default(), mockDataService)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data, _ := mockDataService.ReadOne(1, nil)
	expected, _ := json.Marshal(data)

	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}
