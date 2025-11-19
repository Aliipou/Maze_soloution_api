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

// * This ONLY test that the GetHandler returns the expected response code and body in case of succesfull (200) multiple resource retrieval without the use of a database and the page parameter *
func TestGetHandlerSuccessful(t *testing.T) {
	mockDataService := &service.MockDataServiceSuccessful{}
	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// * GetHanler should call the ReadMany method of the DataService *
	data.GetHandler(rr, req, log.Default(), mockDataService)
	// * Response code should be 200 OK *
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// * We know what the MockDataService will return, so we can compare the response body to the expected value *
	data, _ := mockDataService.ReadMany(0, 10, nil)
	expected, _ := json.Marshal(data)
	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
	}
}

// * This ONLY test that the GetHandler returns the expected response code and body in case of an unsuccesfull (404) multiple resource retrieval without the use of a database and the page parameter *
func TestGetHandlerNotFound(t *testing.T) {
	mockDataService := &service.MockDataServiceNotFound{}
	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// * GetHanler should call the ReadMany method of the DataService *
	data.GetHandler(rr, req, log.Default(), mockDataService)
	// * Response code should be 404 Not Found *
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	// * We know what the MockDataService will return, so we can compare the response body to the expected value *
	if strings.TrimSpace(rr.Body.String()) != `{"error": "Resource not found."}` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `{"error": "Resource not found."}`)
	}
}

// * This ONLY test that the GetHandler returns the expected response code and body in case of an unsuccesfull (500) multiple resource retrieval without the use of a database and the page parameter *
// * Simulates an error in the DataService layer *
func TestGetHandlerError(t *testing.T) {

	mockDataService := &service.MockDataServiceError{}

	req, err := http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	data.GetHandler(rr, req, log.Default(), mockDataService)
	// * Response code should be 500 Internal Server Error *
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	// * We know what the MockDataService will return, so we can compare the response body to the expected value *
	if strings.TrimSpace(rr.Body.String()) != `Internal Server error.` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `Internal Server error.`)
	}
}
