package data_test

import (
	"goapi/internal/api/handlers/data"
	service "goapi/internal/api/service/data"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPutInvalidRequestBody(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Body = io.NopCloser(strings.NewReader(`Plain text, not JSON`))
	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Invalid request data. Please check your input."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerError(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 1, "device_id": "device_id", "device_name": "device_name", "value": 1.0, "type": "type", "date_time": "2020-01-01T00:00:00Z", "description": "description"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceError{})

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Error updating data."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutDataNotFound(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 1, "device_id": "device_id", "device_name": "device_name", "value": 1.0, "type": "type", "date_time": "2020-01-01T00:00:00Z", "description": "description"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceNotFound{})

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error": "Resource not found."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestPutHandlerSuccess(t *testing.T) {

	req, err := http.NewRequest("PUT", "/data", strings.NewReader(`{"id": 1, "device_id": "device_id", "device_name": "device_name", "value": 1.0, "type": "type", "date_time": "2020-01-01T00:00:00Z", "description": "description"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	data.PutHandler(rr, req, log.Default(), &service.MockDataServiceSuccessful{})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"device_id":"device_id","device_name":"device_name","value":1,"type":"type","date_time":"2020-01-01T00:00:00Z","description":"description"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
