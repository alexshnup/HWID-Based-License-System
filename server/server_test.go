package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddKeyHandler(t *testing.T) {
	// Initialize your server, router and etc.

	// Sample request payload
	payload := `{"email": "test@email.com", "expiration": "2023-12-31"}`

	// Create a request to pass to our handler
	req, err := http.NewRequest("POST", "/add", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handler function
	handler := http.HandlerFunc(addKeyHandler)

	// Serve HTTP
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"key added successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
