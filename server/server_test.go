package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var lastTestKey string

func TestAddKeyHandler(t *testing.T) {
	dbfile = "../testdb"
	// Creating a request with a payload
	payload := `{"email": "test@example.com", "expiration": "2023-12-31"}`
	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	// Setting the Authorization header with the token
	req.Header.Set("Authorization", token)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addKeyHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	result := rr.Body.String()

	//lastTestKey={"email":"test@example.com","exp_date":"2023-12-31","license":"8X62-V4AX-OG6Z","message":"New license generated"}
	// parse the json response and get the license key
	var resp map[string]interface{}
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), result)
	}
	lastTestKey = resp["license"].(string)
	fmt.Printf("lastTestKey=%s\n", lastTestKey)

	// Check the response body is what we expect.
	expected := `New license generated`
	// Note that we ignore license details as it is generated randomly
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestListKeysHandler(t *testing.T) {
	dbfile = "../testdb"
	// Creating a request with a payload
	req, err := http.NewRequest("GET", "/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Setting the Authorization header with the token
	req.Header.Set("Authorization", token)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listKeysHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestResetKeyHandler(t *testing.T) {
	dbfile = "../testdb"
	// Creating a request with a payload
	payload := `{"key": "` + lastTestKey + `"}`
	req, err := http.NewRequest("POST", "/reset-key", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	// Setting the Authorization header with the token
	req.Header.Set("Authorization", token)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resetKeyHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
