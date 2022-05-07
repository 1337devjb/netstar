package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"bereths.com/netstar/themoviedb"
)

func TestIndexHandler(t *testing.T) {
	
	request, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(indexHandler)

	hf.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRouter(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	r := newRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "")

	// Handle any unexpected error
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

}

func TestStaticFileServer(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	r := newRouter(themoviedbAPI)
	mockServer := httptest.NewServer(r)

	// get request to assets
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}

func TestSearchHandlerWithInvalidAPIKey(t *testing.T) {
	
	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	request, err := http.NewRequest("GET", "/search?q=Star%20Wars", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(searchHandler(themoviedbAPI))

	hf.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}